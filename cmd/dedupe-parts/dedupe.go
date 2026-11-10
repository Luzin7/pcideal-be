package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/Luzin7/pcideal-be/infra/database"
	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	deleteConfirmation    = "DELETE-DUPLICATES"
	groupOutputLimit      = 10
	indexOptionsConflict  = 85
	indexKeySpecsConflict = 86
)

type duplicateGroup struct {
	Store     string
	URL       string
	Documents []entity.Part
}

type groupKey struct {
	Store string `bson:"store"`
	URL   string `bson:"url"`
}

type duplicateGroupRecord struct {
	Key       groupKey      `bson:"_id"`
	Documents []entity.Part `bson:"documents"`
}

type deletionCounts struct {
	Planned int
	Deleted int
}

type auditSummary struct {
	Groups      int
	MissingKeys int64
}

type verificationExpectation struct {
	CanonicalDocumentIDs []primitive.ObjectID
	PlannedDeletions     int
	Counts               deletionCounts
}

var partFields = []string{"store", "url", "price_cents", "brand", "model", "updated_at"}

type DedupeService struct {
	collection *mongo.Collection
	out        io.Writer
}

func NewDedupeService(collection *mongo.Collection, out io.Writer) *DedupeService {
	return &DedupeService{
		collection: collection,
		out:        out,
	}
}

func (service *DedupeService) DryRun(ctx context.Context) (auditSummary, error) {
	total, err := service.countDocuments(ctx)
	if err != nil {
		return auditSummary{}, err
	}

	missingKeys, err := service.countMissingKeys(ctx)
	if err != nil {
		return auditSummary{}, err
	}

	groups, err := service.auditGroups(ctx)
	if err != nil {
		return auditSummary{}, err
	}

	fmt.Fprintf(service.out, "total parts: %d\n", total)
	fmt.Fprintf(service.out, "duplicate groups: %d\n", len(groups))
	fmt.Fprintf(service.out, "documents in duplicate groups: %d\n", documentsInGroups(groups))
	fmt.Fprintf(service.out, "documents to delete: %d\n", totalDeletions(groups))
	fmt.Fprintf(service.out, "documents with missing store or url: %d\n", missingKeys)
	fmt.Fprintln(service.out, "mode: dry run (no writes)")

	printGroupPlan(service.out, groups)

	return auditSummary{Groups: len(groups), MissingKeys: missingKeys}, nil
}

func (service *DedupeService) Apply(ctx context.Context) error {
	missingKeys, err := service.countMissingKeys(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(service.out, "documents with missing store or url: %d\n", missingKeys)
	if missingKeys > 0 {
		return &MissingKeyError{Count: missingKeys}
	}

	total, err := service.countDocuments(ctx)
	if err != nil {
		return err
	}

	groups, err := service.auditGroups(ctx)
	if err != nil {
		return err
	}

	plannedDeletions := totalDeletions(groups)
	canonicalDocumentIDs := canonicalIDs(groups)

	fmt.Fprintf(service.out, "total parts: %d\n", total)
	fmt.Fprintf(service.out, "duplicate groups: %d\n", len(groups))
	fmt.Fprintf(service.out, "documents in duplicate groups: %d\n", documentsInGroups(groups))
	fmt.Fprintf(service.out, "documents to delete: %d\n", plannedDeletions)

	counts, processed, deleteErr := service.deleteDuplicates(ctx, groups)

	fmt.Fprintf(service.out, "groups processed: %d\n", processed)
	fmt.Fprintf(service.out, "deletions planned at delete time: %d\n", counts.Planned)
	fmt.Fprintf(service.out, "documents deleted: %d\n", counts.Deleted)

	if deleteErr != nil {
		return deleteErr
	}

	if err := service.ensureUniqueIndex(ctx); err != nil {
		return err
	}

	return service.verify(ctx, verificationExpectation{
		CanonicalDocumentIDs: canonicalDocumentIDs,
		PlannedDeletions:     plannedDeletions,
		Counts:               counts,
	})
}

func (service *DedupeService) countDocuments(ctx context.Context) (int64, error) {
	total, err := service.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("count parts: %w", err)
	}

	return total, nil
}

func (service *DedupeService) countMissingKeys(ctx context.Context) (int64, error) {
	filter := bson.M{"$or": []bson.M{
		{"store": bson.M{"$in": []any{"", nil}}},
		{"url": bson.M{"$in": []any{"", nil}}},
	}}

	count, err := service.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("count parts with missing store or url: %w", err)
	}

	return count, nil
}

func (service *DedupeService) auditGroups(ctx context.Context) ([]duplicateGroup, error) {
	cursor, err := service.collection.Aggregate(
		ctx,
		duplicateGroupsPipeline(),
		options.Aggregate().SetAllowDiskUse(true),
	)
	if err != nil {
		return nil, fmt.Errorf("audit duplicate groups: %w", err)
	}
	defer cursor.Close(ctx)

	records := make([]duplicateGroupRecord, 0)
	if err := cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("decode duplicate groups: %w", err)
	}

	groups := make([]duplicateGroup, 0, len(records))
	for _, record := range records {
		groups = append(groups, duplicateGroup{
			Store:     record.Key.Store,
			URL:       record.Key.URL,
			Documents: record.Documents,
		})
	}

	return groups, nil
}

func (service *DedupeService) deleteDuplicates(ctx context.Context, groups []duplicateGroup) (deletionCounts, int, error) {
	counts := deletionCounts{}

	for index, group := range groups {
		groupCounts, err := service.deleteGroup(ctx, group)
		counts.Planned += groupCounts.Planned
		counts.Deleted += groupCounts.Deleted
		if err != nil {
			return counts, index, err
		}

		if index < groupOutputLimit {
			fmt.Fprintf(
				service.out,
				"group %d/%d store=%s url=%s planned=%d deleted=%d\n",
				index+1, len(groups), group.Store, group.URL, groupCounts.Planned, groupCounts.Deleted,
			)
		}
	}

	if len(groups) > groupOutputLimit {
		fmt.Fprintf(service.out, "remaining %d groups processed without per-group output\n", len(groups)-groupOutputLimit)
	}

	return counts, len(groups), nil
}

func (service *DedupeService) deleteGroup(ctx context.Context, group duplicateGroup) (deletionCounts, error) {
	documents, err := service.findGroup(ctx, group.Store, group.URL)
	if err != nil {
		return deletionCounts{}, err
	}

	ids := planDeletions(documents)
	if len(ids) == 0 {
		return deletionCounts{}, nil
	}

	counts := deletionCounts{Planned: len(ids)}

	result, err := service.collection.DeleteMany(ctx, duplicateDeleteFilter(group, ids))
	if err != nil {
		return counts, &DeleteGroupError{Store: group.Store, URL: group.URL, Err: err}
	}

	counts.Deleted = int(result.DeletedCount)

	return counts, nil
}

func (service *DedupeService) findGroup(ctx context.Context, store string, url string) ([]entity.Part, error) {
	cursor, err := service.collection.Find(
		ctx,
		bson.M{"store": store, "url": url},
		options.Find().SetProjection(partProjection()),
	)
	if err != nil {
		return nil, fmt.Errorf("read duplicate group: %w", err)
	}
	defer cursor.Close(ctx)

	documents := make([]entity.Part, 0)
	if err := cursor.All(ctx, &documents); err != nil {
		return nil, fmt.Errorf("decode duplicate group: %w", err)
	}

	return documents, nil
}

func (service *DedupeService) ensureUniqueIndex(ctx context.Context) error {
	err := database.EnsureStoreURLUniqueIndex(ctx, service.collection.Database())
	if err == nil {
		fmt.Fprintf(service.out, "unique index %s: ensured\n", database.StoreURLIndexName)
		return nil
	}

	if !isIndexConflict(err) {
		return err
	}

	fmt.Fprintf(service.out, "unique index %s: conflicting definition, verifying existing index\n", database.StoreURLIndexName)

	return nil
}

func (service *DedupeService) verify(ctx context.Context, expectation verificationExpectation) error {
	reasons := make([]string, 0, 5)

	groups, err := service.auditGroups(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(service.out, "duplicate groups remaining: %d\n", len(groups))
	if len(groups) != 0 {
		reasons = append(reasons, fmt.Sprintf("%d duplicate groups remain", len(groups)))
	}

	indexReady, err := database.UniqueStoreURLIndexReady(ctx, service.collection.Database())
	if err != nil {
		return err
	}

	fmt.Fprintf(service.out, "unique index %s ready: %t\n", database.StoreURLIndexName, indexReady)
	if !indexReady {
		reasons = append(reasons, fmt.Sprintf("index %s is not ready", database.StoreURLIndexName))
	}

	retained, err := service.countCanonicalDocuments(ctx, expectation.CanonicalDocumentIDs)
	if err != nil {
		return err
	}

	fmt.Fprintf(
		service.out,
		"canonical documents retained: %d/%d\n",
		retained, len(expectation.CanonicalDocumentIDs),
	)
	if !canonicalRetentionComplete(len(expectation.CanonicalDocumentIDs), retained) {
		reasons = append(
			reasons,
			fmt.Sprintf("expected %d canonical documents, found %d", len(expectation.CanonicalDocumentIDs), retained),
		)
	}

	fmt.Fprintf(service.out, "deleted documents: %d/%d planned\n", expectation.Counts.Deleted, expectation.PlannedDeletions)
	if expectation.Counts.Deleted != expectation.PlannedDeletions {
		reasons = append(
			reasons,
			fmt.Sprintf("planned %d deletions, deleted %d", expectation.PlannedDeletions, expectation.Counts.Deleted),
		)
	}
	if expectation.Counts.Planned != expectation.PlannedDeletions {
		reasons = append(
			reasons,
			fmt.Sprintf("planned %d deletions, replanned %d at delete time", expectation.PlannedDeletions, expectation.Counts.Planned),
		)
	}

	count, err := service.countDocuments(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(service.out, "final parts count: %d\n", count)

	if len(reasons) != 0 {
		return &VerificationError{Reasons: reasons}
	}

	fmt.Fprintln(service.out, "verification: ok")

	return nil
}

func (service *DedupeService) countCanonicalDocuments(ctx context.Context, ids []primitive.ObjectID) (int64, error) {
	count, err := service.collection.CountDocuments(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return 0, fmt.Errorf("count canonical documents: %w", err)
	}

	return count, nil
}

func duplicateGroupsPipeline() mongo.Pipeline {
	return mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "store", Value: "$store"},
				{Key: "url", Value: "$url"},
			}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "documents", Value: bson.D{{Key: "$push", Value: auditedDocumentPush()}}},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "count", Value: bson.D{{Key: "$gt", Value: 1}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "_id.store", Value: 1},
			{Key: "_id.url", Value: 1},
		}}},
	}
}

func duplicateDeleteFilter(group duplicateGroup, ids []primitive.ObjectID) bson.M {
	return bson.M{
		"_id":   bson.M{"$in": ids},
		"store": group.Store,
		"url":   group.URL,
	}
}

func partProjection() bson.M {
	projection := make(bson.M, len(partFields))

	for _, field := range partFields {
		projection[field] = 1
	}

	return projection
}

func auditedDocumentPush() bson.D {
	document := bson.D{{Key: "_id", Value: "$_id"}}

	for _, field := range partFields {
		document = append(document, bson.E{Key: field, Value: "$" + field})
	}

	return document
}

func isIndexConflict(err error) bool {
	var commandError mongo.CommandError
	if !errors.As(err, &commandError) {
		return false
	}

	return commandError.Code == indexOptionsConflict || commandError.Code == indexKeySpecsConflict
}

func pickCanonical(documents []entity.Part) entity.Part {
	if len(documents) == 0 {
		return entity.Part{}
	}

	canonical := documents[0]
	for _, document := range documents[1:] {
		if isNewerDocument(document, canonical) {
			canonical = document
		}
	}

	return canonical
}

func isNewerDocument(candidate entity.Part, current entity.Part) bool {
	if candidate.UpdatedAt.Equal(current.UpdatedAt) {
		return bytes.Compare(candidate.ID[:], current.ID[:]) > 0
	}

	return candidate.UpdatedAt.After(current.UpdatedAt)
}

func planDeletions(documents []entity.Part) []primitive.ObjectID {
	if len(documents) < 2 {
		return nil
	}

	canonicalID := pickCanonical(documents).ID
	ids := make([]primitive.ObjectID, 0, len(documents)-1)

	for _, document := range documents {
		if document.ID != canonicalID {
			ids = append(ids, document.ID)
		}
	}

	sort.Slice(ids, func(i int, j int) bool {
		return bytes.Compare(ids[i][:], ids[j][:]) < 0
	})

	return ids
}

func canonicalIDs(groups []duplicateGroup) []primitive.ObjectID {
	ids := make([]primitive.ObjectID, 0, len(groups))

	for _, group := range groups {
		if len(group.Documents) == 0 {
			continue
		}

		ids = append(ids, pickCanonical(group.Documents).ID)
	}

	return ids
}

func canonicalRetentionComplete(expected int, retained int64) bool {
	return retained == int64(expected)
}

func differingFields(documents []entity.Part) []string {
	if len(documents) < 2 {
		return nil
	}

	fields := make([]string, 0, 4)
	if hasPriceDifference(documents) {
		fields = append(fields, "price_cents")
	}
	if hasBrandDifference(documents) {
		fields = append(fields, "brand")
	}
	if hasModelDifference(documents) {
		fields = append(fields, "model")
	}
	if hasUpdatedAtDifference(documents) {
		fields = append(fields, "updated_at")
	}

	return fields
}

func hasPriceDifference(documents []entity.Part) bool {
	return hasDifference(documents, func(document entity.Part) int64 {
		return document.PriceCents
	})
}

func hasBrandDifference(documents []entity.Part) bool {
	return hasDifference(documents, func(document entity.Part) string {
		return document.Brand
	})
}

func hasModelDifference(documents []entity.Part) bool {
	return hasDifference(documents, func(document entity.Part) string {
		return document.Model
	})
}

func hasUpdatedAtDifference(documents []entity.Part) bool {
	return hasDifference(documents, func(document entity.Part) int64 {
		return document.UpdatedAt.UnixNano()
	})
}

func hasDifference[T comparable](documents []entity.Part, extract func(entity.Part) T) bool {
	if len(documents) < 2 {
		return false
	}

	first := extract(documents[0])

	for _, document := range documents[1:] {
		if extract(document) != first {
			return true
		}
	}

	return false
}

func totalDeletions(groups []duplicateGroup) int {
	total := 0
	for _, group := range groups {
		total += deletionsInGroup(group)
	}

	return total
}

func documentsInGroups(groups []duplicateGroup) int {
	total := 0
	for _, group := range groups {
		total += len(group.Documents)
	}

	return total
}

func deletionsInGroup(group duplicateGroup) int {
	if len(group.Documents) < 2 {
		return 0
	}

	return len(group.Documents) - 1
}

func printGroupPlan(out io.Writer, groups []duplicateGroup) {
	limit := len(groups)
	if limit > groupOutputLimit {
		limit = groupOutputLimit
	}

	for index := 0; index < limit; index++ {
		group := groups[index]
		fmt.Fprintf(
			out,
			"group %d/%d store=%s url=%s deletes=%d differs=%s\n",
			index+1, len(groups), group.Store, group.URL,
			deletionsInGroup(group), formatDifferingFields(differingFields(group.Documents)),
		)
	}

	if len(groups) > groupOutputLimit {
		fmt.Fprintf(out, "showing first %d of %d groups\n", groupOutputLimit, len(groups))
	}
}

func formatDifferingFields(fields []string) string {
	if len(fields) == 0 {
		return "none"
	}

	return strings.Join(fields, ",")
}
