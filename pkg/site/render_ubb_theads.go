package site

import (
	"github.com/tivizi/forarun/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ThreadsRender 主题列表渲染器
type ThreadsRender struct {
	BaseRender
}

// Render 渲染UUB
func (t *ThreadsRender) Render(raw string) string {
	arg0, err := t.strArg(0)
	if err != nil {
		return "«expr error: " + err.Error() + "»"
	}
	size, err := t.int64Arg(1)
	if err != nil {
		return "«expr error: " + err.Error() + "»"
	}
	bbsID, err := t.strArg(2)
	if err != nil {
		return "«expr error: " + err.Error() + "»"
	}
	var options options.FindOptions
	var threads []*domain.Thread
	switch arg0 {
	case "new":
		options.SetSort(bson.M{"createtime": -1})
	case "active":
		options.SetSort(bson.M{"lastactivetime": -1})
	default:
		return "«unsupported threads type»"
	}
	if bbsID != "0" {
		objID, err := primitive.ObjectIDFromHex(bbsID)
		if err != nil {
			return "«expr error: " + err.Error() + "»"
		}
		threads, err = domain.LoadThreadsByFilterAndOptions(t.Site.ID, size, &bson.M{
			"bbscontexts": bson.M{"$elemMatch": bson.M{
				"_id": objID,
			}},
		}, &options)
	} else {
		threads, err = domain.LoadThreadsByFilterAndOptions(t.Site.ID, size, &bson.M{}, &options)
	}
	if err != nil {
		return "«load error: " + err.Error() + "»"
	}
	return t.renderWithTemplate("segment_threads_list.html", threads)
}
