package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	dbinfra "mockhu-app-backend/internal/infra/db"
)

func main() {
	ctx := context.Background()

	// Connect to database
	pg, err := dbinfra.New(ctx, dbinfra.DatabaseURLFromEnv())
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}
	defer pg.Close()

	log.Println("✅ Database connected")

	// Get existing users
	rows, err := pg.Pool.Query(ctx, "SELECT id, email, username FROM users ORDER BY created_at LIMIT 10")
	if err != nil {
		log.Fatalf("Failed to query users: %v", err)
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var id, email string
		var username *string
		if err := rows.Scan(&id, &email, &username); err != nil {
			log.Printf("Failed to scan user: %v", err)
			continue
		}
		userIDs = append(userIDs, id)
		usernameStr := "N/A"
		if username != nil {
			usernameStr = *username
		}
		log.Printf("Found user: %s (%s) - username: %s", email, id, usernameStr)
	}

	if len(userIDs) == 0 {
		log.Fatal("No users found. Please create users first before seeding posts.")
	}

	log.Printf("Found %d users. Creating posts...\n", len(userIDs))

	// Sample post contents
	postContents := []struct {
		content     string
		images      []string
		isAnonymous bool
	}{
		{"Just finished building my first API with Go and Fiber! The performance is incredible. 🚀 #coding #golang", []string{"https://images.unsplash.com/photo-1516116216624-53e697fedbea?w=800"}, false},
		{"Working on a new social media platform. The architecture is coming together nicely. Excited to share more soon! 💻", []string{}, false},
		{"Anyone else love the feeling when your code finally compiles after hours of debugging? 😅", []string{}, false},
		{"Just discovered PostgreSQL arrays. Mind blown! 🤯 The flexibility is amazing.", []string{"https://images.unsplash.com/photo-1544383835-b9af0e3b90f9?w=800"}, false},
		{"Beautiful sunset today! 🌅 Sometimes you just need to stop and appreciate the little things.", []string{"https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=800"}, false},
		{"Started learning photography. Here's my first attempt at street photography. Feedback welcome! 📷", []string{"https://images.unsplash.com/photo-1492684223066-81342ee5ff30?w=800"}, false},
		{"Coffee and coding. The perfect morning routine ☕", []string{}, false},
		{"Just finished reading an amazing book on design systems. The principles apply everywhere! 📚", []string{}, false},
		{"Working out at the gym today. Consistency is key! 💪", []string{}, false},
		{"New recipe tried today: Homemade pasta. Turned out amazing! 🍝", []string{"https://images.unsplash.com/photo-1551462147-5bc923f49fef?w=800"}, false},
		{"Travel tip: Always pack a power bank. You'll thank yourself later! ✈️", []string{}, false},
		{"Anonymous post: Sometimes I wonder if anyone actually reads these... 🤔", []string{}, true},
		{"Weekend project: Building a REST API with proper error handling. Learning so much! 📖", []string{}, false},
		{"Code review tip: Always explain the 'why' not just the 'what'. Makes a huge difference! 💡", []string{}, false},
		{"Just deployed my first production API. The feeling is unreal! 🎉", []string{"https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=800"}, false},
	}

	rand.Seed(time.Now().UnixNano())

	// Create posts for each user
	postCount := 0
	for i, userID := range userIDs {
		// Create 3-5 posts per user
		numPosts := rand.Intn(3) + 3
		for j := 0; j < numPosts && j < len(postContents); j++ {
			content := postContents[(i*numPosts+j)%len(postContents)]

			// Randomize view count
			viewCount := rand.Intn(200)

			// Randomize creation time (within last 3 days)
			createdAt := time.Now().Add(-time.Duration(rand.Intn(72)) * time.Hour)

			query := `
				INSERT INTO posts (user_id, content, images, is_anonymous, view_count, created_at)
				VALUES ($1, $2, $3, $4, $5, $6)
				RETURNING id
			`

			var postID string
			err := pg.Pool.QueryRow(ctx, query, userID, content.content, content.images, content.isAnonymous, viewCount, createdAt).
				Scan(&postID)
			if err != nil {
				log.Printf("Failed to create post: %v", err)
				continue
			}

			postCount++
			log.Printf("Created post %s for user %s", postID, userID)

			// Add some reactions (30% chance)
			if rand.Float32() < 0.3 {
				// Get random users to react
				numReactions := rand.Intn(5) + 1
				for k := 0; k < numReactions && k < len(userIDs); k++ {
					reactorID := userIDs[rand.Intn(len(userIDs))]
					if reactorID != userID { // Don't react to own post
						reactionQuery := `
							INSERT INTO post_reactions (post_id, user_id, reaction_type, created_at)
							VALUES ($1, $2, 'fire', $3)
							ON CONFLICT (post_id, user_id) DO NOTHING
						`
						_, err := pg.Pool.Exec(ctx, reactionQuery, postID, reactorID, createdAt.Add(time.Duration(rand.Intn(60))*time.Minute))
						if err != nil {
							log.Printf("Failed to add reaction: %v", err)
						}
					}
				}
			}
		}
	}

	log.Printf("\n✅ Successfully created %d posts with reactions!", postCount)
}
