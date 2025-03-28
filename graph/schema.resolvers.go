package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.68

import (
	"context"
	"errors"
	"graphql-crud/auth"
	"graphql-crud/graph/model"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, userID string, name *string, email *string) (*model.User, error) {
	collection := r.Client.Database(os.Getenv("DB_NAME")).Collection("users")
	userCheck, _ := r.Query().User(ctx, userID)
	if userCheck == nil {
		return nil, errors.New("User not found")
	}

	var userDuplicate model.User
	_ = collection.FindOne(ctx, bson.M{"email": email}).Decode(&userDuplicate)
	if userDuplicate.Email != "" && userDuplicate.Email != userCheck.Email {
		return nil, errors.New("email already registered")
	}

	user := &model.User{
		Name:  *name,
		Email: *email,
	}

	objID, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"name":  user.Name,
			"email": user.Email,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, userID string) (*model.User, error) {
	collection := r.Client.Database(os.Getenv("DB_NAME")).Collection("users")
	user, _ := r.Query().User(ctx, userID)
	if user == nil {
		return nil, errors.New("User not found")
	}

	objID, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{"_id": objID}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// CreateProduct is the resolver for the createProduct field.
func (r *mutationResolver) CreateProduct(ctx context.Context, name string, price int32, stock int32) (*model.Product, error) {
	collection := r.Client.Database(os.Getenv("DB_NAME")).Collection("products")
	product := model.Product{
		Name:  name,
		Price: price,
		Stock: stock,
	}

	_, err := collection.InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// UpdateProduct is the resolver for the updateProduct field.
func (r *mutationResolver) UpdateProduct(ctx context.Context, productID string, name *string, price *int32, stock *int32) (*model.Product, error) {
	collection := r.Client.Database(os.Getenv("DB_NAME")).Collection("products")
	productCheck, _ := r.Query().Product(ctx, productID)
	if productCheck == nil {
		return nil, errors.New("Product not found")
	}

	product := model.Product{
		Name:  *name,
		Price: *price,
		Stock: *stock,
	}

	objID, _ := primitive.ObjectIDFromHex(productID)
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"name":  product.Name,
			"price": product.Price,
			"stock": product.Stock,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// DeleteProduct is the resolver for the deleteProduct field.
func (r *mutationResolver) DeleteProduct(ctx context.Context, productID string) (*model.Product, error) {
	collection := r.Client.Database(os.Getenv("DB_NAME")).Collection("products")
	product, err := r.Query().Product(ctx, productID)
	if err != nil {
		return nil, errors.New("Product not found")
	}

	objID, _ := primitive.ObjectIDFromHex(productID)
	filter := bson.M{"_id": objID}
	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.AuthPayload, error) {
	var userParse model.UserParse
	collection := r.Client.Database(os.Getenv("DB_NAME")).Collection("users")
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&userParse)
	if err != nil {
		return nil, errors.New("invalid email/password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userParse.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email/password")
	}

	token, err := auth.GenerateToken(userParse.ID.Hex())
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	user := model.User{
		Name:   userParse.Name,
		UserID: userParse.ID.Hex(),
		Email:  userParse.Email,
	}

	return &model.AuthPayload{
		Token: token,
		User:  &user,
	}, nil
}

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, name string, email string, password string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	collection := r.Client.Database(os.Getenv("DB_NAME")).Collection("users")
	var userCheck model.User
	_ = collection.FindOne(ctx, bson.M{"email": email}).Decode(&userCheck)
	if userCheck.Email != "" {
		return nil, errors.New("email already registered")
	}

	user := model.User{Name: name, Email: email, Password: string(hashedPassword)}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	collection := r.Client.Database(os.Getenv("DB_NAME")).Collection("users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var users []*model.User
	for cursor.Next(ctx) {
		var userParse *model.UserParse
		if err := cursor.Decode(&userParse); err != nil {
			return nil, err
		}

		userAppend := model.User{
			Name:   userParse.Name,
			UserID: userParse.ID.Hex(),
			Email:  userParse.Email,
		}

		users = append(users, &userAppend)
	}

	return users, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, userID string) (*model.User, error) {
	var user model.User

	collection := r.Client.Database(os.Getenv("DB_NAME")).Collection("users")

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	user.UserID = userID

	return &user, nil
}

// Products is the resolver for the products field.
func (r *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {
	collection := r.Client.Database(os.Getenv("DB_NAME")).Collection("products")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var products []*model.Product
	for cursor.Next(ctx) {
		var productParse *model.ProductParse
		if err := cursor.Decode(&productParse); err != nil {
			return nil, err
		}

		product := model.Product{
			Name:      productParse.Name,
			Stock:     productParse.Stock,
			Price:     productParse.Price,
			ProductID: productParse.ID.Hex(),
		}

		products = append(products, &product)
	}

	return products, nil
}

// Product is the resolver for the product field.
func (r *queryResolver) Product(ctx context.Context, productID string) (*model.Product, error) {
	collection := r.Client.Database(os.Getenv("DB_NAME")).Collection("products")

	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return nil, err
	}

	var product model.Product
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		return nil, err
	}

	product.ProductID = productID
	return &product, nil
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	var user model.User
	collection := r.Client.Database(os.Getenv("DB_NAME")).Collection("users")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
