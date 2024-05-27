package repository

import (
    "context"
    "go-microservice-demo/pkg/model"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository struct {
    collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
    return &ProductRepository{
        collection: db.Collection("product"),
    }
}

func (r *ProductRepository) Save(product *model.Product) error {
    _, err := r.collection.InsertOne(context.Background(), product)
    return err
}

func (r *ProductRepository) FindAll() ([]model.Product, error) {
    var products []model.Product
    cursor, err := r.collection.Find(context.Background(), bson.D{{}}, options.Find())
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())
    for cursor.Next(context.Background()) {
        var product model.Product
        if err := cursor.Decode(&product); err != nil {
            return nil, err
        }
        products = append(products, product)
    }
    return products, nil
}
