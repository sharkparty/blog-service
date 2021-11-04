package server

import (
	"context"
	"fmt"

	config "blog-service/config"
	blogProto "blog-service/rpc/blog"

	"github.com/twitchtv/twirp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct{}

type BlogItem struct {
	Id      primitive.ObjectID `bson:"_id"`
	Title   string             `bson:"title"`
	Content string             `bson:"content"`
}

func (*Server) CreateBlog(ctc context.Context, req *blogProto.CreateBlogRequest) (*blogProto.CreateBlogResponse, error) {
	data := &blogProto.CreateBlogRequest{
		Title:   req.GetTitle(),
		Content: req.GetContent(),
	}

	res, err := config.Collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("There was an error creating a blog: %v", err))
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("Cannot convert to oid: %v", ok))
	}

	return &blogProto.CreateBlogResponse{
		Id:      oid.Hex(),
		Title:   data.Title,
		Content: data.Content,
	}, nil
}

func (*Server) GetBlog(ctc context.Context, req *blogProto.GetBlogRequest) (*blogProto.GetBlogResponse, error) {
	id := req.GetId()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "Invalid blog ID")
	}
	filter := bson.D{{Key: "_id", Value: oid}}

	result := BlogItem{}
	merr := config.Collection.FindOne(context.TODO(), filter).Decode(&result)
	if merr != nil {
		if merr == mongo.ErrNoDocuments {
			return nil, twirp.NewError(twirp.NotFound, fmt.Sprintf("No documents were found for id: %v", id))
		}
		return nil, twirp.NewError(twirp.NotFound, fmt.Sprintf("There was an error finding a blog with ID: %v \nError: %v", id, merr))
	}

	return &blogProto.GetBlogResponse{
		Id:      id,
		Title:   result.Title,
		Content: result.Content,
	}, nil
}

func (*Server) UpdateBlog(ctc context.Context, req *blogProto.UpdateBlogRequest) (*blogProto.UpdateBlogResponse, error) {
	id := req.GetId()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "Invalid blog ID")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: oid}}
	newBlog := &blogProto.Blog{
		Title:   req.GetTitle(),
		Content: req.GetContent(),
	}
	update := bson.D{{"$set", bson.M{"title": newBlog.Title, "content": newBlog.Content}}}

	cursor, merr := config.Collection.UpdateOne(context.TODO(), filter, update)
	if merr != nil || cursor.MatchedCount == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("Blog id: %v could not be updated with %v", id, newBlog))
	}

	return &blogProto.UpdateBlogResponse{
		Id:      oid.Hex(),
		Title:   req.GetTitle(),
		Content: req.GetContent(),
	}, nil
}

func (*Server) DeleteBlog(ctc context.Context, req *blogProto.DeleteBlogRequest) (*blogProto.DeleteBlogResponse, error) {
	id := req.GetId()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "Invalid blog ID")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: oid}}
	cursor, err := config.Collection.DeleteOne(context.TODO(), filter)

	if err != nil || cursor.DeletedCount != 1 {
		return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("Unable to delete blog with ID: %v", id))
	}
	return &blogProto.DeleteBlogResponse{
		Id: id,
	}, nil
}

func (*Server) ListBlog(ctc context.Context, req *blogProto.ListBlogRequest) (*blogProto.ListBlogResponse, error) {
	filter := bson.D{}
	limit := int64(25)

	if req.GetLimit() > 0 {
		limit = req.GetLimit()
	}

	options := &options.FindOptions{
		Limit: &limit,
	}

	var results []BlogItem
	cursor, merr := config.Collection.Find(context.TODO(), filter, options)

	if merr != nil {
		if merr == mongo.ErrNoDocuments {
			return nil, twirp.NewError(twirp.NotFound, "No documents were found")
		}
		return nil, twirp.NewError(twirp.NotFound, "There was an error listing blogs")
	}

	if err := cursor.All(context.TODO(), &results); err != nil {
		fmt.Printf("error: %v", err)
	}

	blogs := []*blogProto.CreateBlogResponse{}

	for _, result := range results {
		blog := blogProto.CreateBlogResponse{
			Id:      result.Id.Hex(),
			Title:   result.Title,
			Content: result.Content,
		}
		blogs = append(blogs, &blog)
	}

	return &blogProto.ListBlogResponse{
		Blogs: blogs,
	}, nil
}
