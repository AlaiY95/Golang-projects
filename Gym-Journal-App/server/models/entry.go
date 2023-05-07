package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Define a struct called Entry that represents a workout log entry.
// All fields are pointers to allow them to be optional.
type Entry struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Exercises *string            `bson:"exercises,omitempty" json:"exercises,omitempty"`
	Muscle    *string            `bson:"muscle,omitempty" json:"muscle,omitempty"`
	Sets      *int               `bson:"sets,omitempty" json:"sets,omitempty"`
	Reps      *int               `bson:"reps,omitempty" json:"reps,omitempty"`
}
