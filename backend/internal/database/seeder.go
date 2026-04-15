package database

import (
	"context"
	"log"
	"time"

	"tll-backend/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedDatabase initializes the database with sample data if it's empty
func SeedDatabase(db *gorm.DB) error {
	ctx := context.Background()

	// Count existing users
	var userCount int64
	if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		log.Printf("Warning: Could not count users for seeding check: %v", err)
		return nil // Don't fail on seeding
	}

	// Only seed if database is empty
	// Ignore system user that got created during migration, so check for more than 1 user
	if userCount > 1 {
		log.Println("Database already has users, skipping seeding")
		return nil
	}

	log.Println("Database is empty, seeding with sample data...")

	// Sample users to seed
	sampleUsers := []struct {
		email       string
		username    string
		password    string
		displayName string
		role        string
		bio         string
	}{
		{
			email:       "admin@travellink.local",
			username:    "admin",
			password:    "AdminPass123!",
			displayName: "Administrator",
			role:        models.RoleAdmin.String(),
			bio:         "System administrator",
		},
		{
			email:       "traveller@travellink.local",
			username:    "traveller",
			password:    "TravellerPass123!",
			displayName: "Travel Enthusiast",
			role:        models.RoleTraveller.String(),
			bio:         "Loves exploring new places",
		},
		{
			email:       "user@travellink.local",
			username:    "user",
			password:    "UserPass123!",
			displayName: "Simple User",
			role:        models.RoleSimple.String(),
			bio:         "Just browsing travel plans",
		},
		{
			email:       "alice@travellink.local",
			username:    "alice",
			password:    "AlicePass123!",
			displayName: "Alice Cooper",
			role:        models.RoleTraveller.String(),
			bio:         "Adventure seeker and photographer",
		},
		{
			email:       "bob@travellink.local",
			username:    "bob",
			password:    "BobPass123!",
			displayName: "Bob Smith",
			role:        models.RoleTraveller.String(),
			bio:         "Cultural explorer",
		},
		{
			email:       "charlie@travellink.local",
			username:    "charlie",
			password:    "CharliePass123!",
			displayName: "Charlie Brown",
			role:        models.RoleSimple.String(),
			bio:         "Beach and nature lover",
		},
	}

	// Create users and store their IDs for plan creation
	userMap := make(map[string]string) // username -> userID
	for _, userData := range sampleUsers {
		// Hash password using bcrypt with default cost
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password for %s: %v", userData.username, err)
			continue
		}

		// Create user model
		userID := uuid.New().String()
		user := &models.User{
			ID:           userID,
			Email:        userData.email,
			Username:     userData.username,
			PasswordHash: string(hashedPassword),
			DisplayName:  userData.displayName,
			Bio:          userData.bio,
			Role:         userData.role,
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		// Create user directly in database via GORM
		if err := db.WithContext(ctx).Create(user).Error; err != nil {
			log.Printf("Warning: Could not seed user %s: %v", userData.username, err)
			continue
		}

		userMap[userData.username] = userID
		log.Printf("✓ Seeded user: %s (%s)", userData.username, userData.email)
	}

	// Sample travel plans to seed
	samplePlans := []struct {
		title       string
		description string
		destination string
		authorName  string
		status      string
	}{
		{
			title:       "Summer European Adventure",
			description: "Experience the best of Europe during summer with this comprehensive 3-week itinerary covering Paris, Rome, and Barcelona.",
			destination: "Europe",
			authorName:  "alice",
			status:      "published",
		},
		{
			title:       "Tokyo Food Tour",
			description: "A culinary journey through Tokyo exploring traditional and modern Japanese cuisine.",
			destination: "Tokyo, Japan",
			authorName:  "bob",
			status:      "published",
		},
		{
			title:       "New Zealand Road Trip",
			description: "Epic road trip around both islands of New Zealand with stunning scenery and outdoor activities.",
			destination: "New Zealand",
			authorName:  "alice",
			status:      "published",
		},
		{
			title:       "Budget Southeast Asia Backpacking",
			description: "Affordable backpacking route through Thailand, Vietnam, and Cambodia.",
			destination: "Southeast Asia",
			authorName:  "charlie",
			status:      "draft",
		},
		{
			title:       "Caribbean Island Hopping",
			description: "Visit multiple Caribbean islands for beaches, diving, and tropical paradise experiences.",
			destination: "Caribbean",
			authorName:  "bob",
			status:      "draft",
		},
		{
			title:       "Machu Picchu and Peruvian Highlands",
			description: "Trek through the breathtaking Peruvian Andes and explore the ancient Inca ruins.",
			destination: "Peru",
			authorName:  "traveller",
			status:      "published",
		},
		{
			title:       "Norway Fjords Road Trip",
			description: "Scenic drive through Norwegian fjords with hiking and photography opportunities.",
			destination: "Norway",
			authorName:  "alice",
			status:      "draft",
		},
		{
			title:       "Morocco Cultural Immersion",
			description: "Explore the markets, deserts, and mountain villages of Morocco.",
			destination: "Morocco",
			authorName:  "bob",
			status:      "published",
		},
		{
			title:       "Iceland Ring Road Adventure",
			description: "Complete circle around Iceland exploring waterfalls, glaciers, and geysers.",
			destination: "Iceland",
			authorName:  "charlie",
			status:      "published",
		},
		{
			title:       "India Taj Mahal and Beyond",
			description: "Explore India's cultural heritage from the Taj Mahal to Kerala backwaters.",
			destination: "India",
			authorName:  "traveller",
			status:      "draft",
		},
	}

	// Create travel plans and their associated nodes
	planMap := make(map[string]string) // planTitle -> planID

	for _, planData := range samplePlans {
		authorID, exists := userMap[planData.authorName]
		if !exists {
			log.Printf("Warning: Author %s not found, skipping plan creation", planData.authorName)
			continue
		}

		plan := &models.TravelPlan{
			ID:           uuid.New().String(),
			Title:        planData.title,
			Description:  planData.description,
			Destination:  planData.destination,
			AuthorID:     authorID,
			Status:       planData.status,
			RatingCount:  0,
			RatingSum:    0,
			CommentCount: 0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := db.WithContext(ctx).Create(plan).Error; err != nil {
			log.Printf("Warning: Could not seed plan %s: %v", planData.title, err)
			continue
		}

		planMap[planData.title] = plan.ID
		log.Printf("✓ Seeded plan: %s (status: %s, author: %s)", planData.title, planData.status, planData.authorName)
	}

	// STEP 1: Create generic reusable nodes (independent of any plan)
	// These nodes are created once and reused across plans
	adminID := userMap["admin"]
	nodeLibrary := make(map[string]string) // nodeKey -> nodeID

	// Generic attraction node definitions (reusable across plans)
	attractionDefs := []struct {
		key      string // unique identifier for node library
		name     string
		category string
		location string
		desc     string
	}{
		// Paris attractions
		{key: "paris-eiffel", name: "Eiffel Tower", category: "tourist_attraction", location: "Champ de Mars, Paris, France", desc: "Iconic Iron Lady offering panoramic city views"},
		{key: "paris-louvre", name: "Louvre Museum", category: "museum", location: "Palais-Royal, Paris, France", desc: "World's largest art museum with masterpieces"},
		{key: "paris-notredame", name: "Notre-Dame", category: "tourist_attraction", location: "Île de la Cité, Paris, France", desc: "Historic Gothic cathedral"},

		// Rome attractions
		{key: "rome-colosseum", name: "Colosseum", category: "tourist_attraction", location: "Piazza del Colosseo, Rome, Italy", desc: "Ancient Roman amphitheater"},
		{key: "rome-vatican", name: "Vatican City", category: "tourist_attraction", location: "Vatican City", desc: "World's smallest country with St. Peter's Basilica"},
		{key: "rome-forum", name: "Roman Forum", category: "tourist_attraction", location: "Via della Salara Vecchia, Rome, Italy", desc: "Ancient center of Roman public life"},

		// Barcelona attractions
		{key: "barcelona-sagrada", name: "Sagrada Familia", category: "tourist_attraction", location: "Barcelona, Spain", desc: "Gaudí's stunning unfinished basilica"},
		{key: "barcelona-park", name: "Park Güell", category: "park", location: "Barcelona, Spain", desc: "Modernist park with colorful mosaics"},

		// Tokyo attractions
		{key: "tokyo-tsukiji", name: "Tsukiji Outer Market", category: "restaurant", location: "Tsukiji, Tokyo", desc: "Famous seafood market with street food"},
		{key: "tokyo-shibuya", name: "Shibuya Crossing", category: "tourist_attraction", location: "Shibuya, Tokyo", desc: "World's busiest pedestrian crossing"},
		{key: "tokyo-ramen", name: "Ramen Alley", category: "restaurant", location: "Shinjuku, Tokyo", desc: "Tiny alley with excellent ramen restaurants"},
		{key: "tokyo-meiji", name: "Meiji Shrine", category: "tourist_attraction", location: "Shibuya, Tokyo", desc: "Shinto shrine in peaceful forest"},

		// New Zealand attractions
		{key: "nz-milford", name: "Milford Sound", category: "park", location: "Fiordland, South Island, New Zealand", desc: "UNESCO fjord with waterfalls and wildlife"},
		{key: "nz-queenstown", name: "Queenstown Adventure Sports", category: "park", location: "Queenstown, New Zealand", desc: "Bungy jumping, jet boating, skydiving capital"},
		{key: "nz-glow-worms", name: "Glow Worm Caves", category: "tourist_attraction", location: "Te Anau, New Zealand", desc: "Underground caves with bioluminescent glow worms"},
		{key: "nz-hobbiton", name: "Hobbiton Movie Set", category: "tourist_attraction", location: "Matamata, North Island, New Zealand", desc: "Lord of the Rings filming location"},

		// Southeast Asia attractions
		{key: "bangkok-palace", name: "Grand Palace", category: "tourist_attraction", location: "Bangkok, Thailand", desc: "Thailand's most sacred Buddhist temple complex"},
		{key: "phuket-beach", name: "Phuket Beaches", category: "park", location: "Phuket, Thailand", desc: "Beautiful beaches with clear water"},
		{key: "hanoi-oldquarter", name: "Hanoi Old Quarter", category: "restaurant", location: "Hanoi, Vietnam", desc: "Chaotic historic quarter with street food"},
		{key: "halong-bay", name: "Halong Bay", category: "park", location: "Halong Bay, Vietnam", desc: "UNESCO World Heritage bay with limestone cliffs"},
		{key: "angkor-wat", name: "Angkor Wat", category: "tourist_attraction", location: "Siem Reap, Cambodia", desc: "Ancient Khmer temple complex and UNESCO site"},

		// Peru attractions
		{key: "cusco-city", name: "Cusco Historic Center", category: "tourist_attraction", location: "Cusco, Peru", desc: "Colonial city with Inca heritage"},
		{key: "machu-picchu", name: "Machu Picchu", category: "tourist_attraction", location: "Aguas Calientes, Peru", desc: "Ancient Incan citadel among Andean peaks"},
		{key: "cusco-market", name: "Sacred Valley Markets", category: "shopping", location: "Ollantaytambo, Peru", desc: "Indigenous market with traditional crafts"},

		// Iceland attractions
		{key: "iceland-gullfoss", name: "Gullfoss Waterfall", category: "park", location: "Golden Circle, Iceland", desc: "Massive two-stage waterfall"},
		{key: "iceland-black-sand", name: "Black Sand Beach", category: "park", location: "Reynisfjara, Iceland", desc: "Otherworldly black sand beach with basalt cliffs"},
		{key: "iceland-geysir", name: "Geysir Hot Springs", category: "park", location: "Haukadalur Valley, Iceland", desc: "Active geysers and colorful hot spring pools"},
		{key: "iceland-blue-lagoon", name: "Blue Lagoon", category: "park", location: "Reykjanes Peninsula, Iceland", desc: "Geothermal spa with milky blue mineral water"},

		// Caribbean attractions
		{key: "cayman-beach", name: "Seven Mile Beach", category: "park", location: "Grand Cayman, Caribbean", desc: "World-class white-sand beach"},
		{key: "jamaica-falls", name: "Dunn's River Falls", category: "park", location: "Ocho Rios, Jamaica", desc: "Cascading waterfall flowing into natural pools"},
		{key: "turks-grace", name: "Grace Bay Beach", category: "park", location: "Turks and Caicos", desc: "Pristine Caribbean beach"},
	}

	// Create generic attraction nodes
	for _, adef := range attractionDefs {
		node := &models.Node{
			ID:         uuid.New().String(),
			Type:       "attraction",
			CreatedBy:  adminID,
			IsApproved: true,
			CreatedAt:  time.Now(),
		}
		if err := db.WithContext(ctx).Create(node).Error; err != nil {
			log.Printf("Warning: Could not create attraction node %s: %v", adef.key, err)
			continue
		}

		attraction := &models.AttractionNodeDetail{
			NodeID:           node.ID,
			Name:             adef.name,
			Category:         adef.category,
			Location:         adef.location,
			Description:      adef.desc,
			HoursOfOperation: "09:00-18:00",
			CreatedAt:        time.Now(),
		}
		if err := db.WithContext(ctx).Create(attraction).Error; err != nil {
			log.Printf("Warning: Could not create attraction detail %s: %v", adef.key, err)
			continue
		}

		nodeLibrary[adef.key] = node.ID
	}

	// Generic transition node definitions (reusable across plans)
	transitionDefs := []struct {
		key   string // unique identifier
		title string // e.g., "International Flight", "Subway Line 5"
		mode  string
		desc  string
	}{
		{key: "flight-intl", title: "International Flight", mode: "flight", desc: "Long-haul flight between countries"},
		{key: "flight-regional", title: "Regional Flight", mode: "flight", desc: "Short regional flight between destinations"},
		{key: "train-overnight", title: "Overnight Train", mode: "train", desc: "Comfortable overnight train journey"},
		{key: "train-express", title: "Express Train", mode: "train", desc: "Fast express train service"},
		{key: "walk-short", title: "Walking", mode: "walking", desc: "Short walk between nearby attractions"},
		{key: "drive-scenic", title: "Scenic Drive", mode: "car", desc: "Leisurely road trip between destinations"},
		{key: "drive-highway", title: "Highway Drive", mode: "car", desc: "Interstate highway travel"},
		{key: "subway-tokyo", title: "Tokyo Subway", mode: "bus", desc: "Tokyo metro system trains"},
		{key: "subway-europe", title: "European Metro", mode: "bus", desc: "Urban metro/subway transit"},
		{key: "bus-transit", title: "Public Bus", mode: "bus", desc: "Local and intercity buses"},
		{key: "ferry", title: "Ferry", mode: "bus", desc: "Water transport between islands or across bays"},
		{key: "taxi", title: "Taxi", mode: "taxi", desc: "Point-to-point taxi service"},
		{key: "taxi-shibuya", title: "Taxi to Shibuya", mode: "taxi", desc: "Taxi ride to Shibuya district"},
		{key: "bike", title: "Bike Rental", mode: "bike", desc: "Self-guided cycling between attractions"},
	}

	// Create generic transition nodes
	for _, tdef := range transitionDefs {
		node := &models.Node{
			ID:         uuid.New().String(),
			Type:       "transition",
			CreatedBy:  adminID,
			IsApproved: true,
			CreatedAt:  time.Now(),
		}
		if err := db.WithContext(ctx).Create(node).Error; err != nil {
			log.Printf("Warning: Could not create transition node %s: %v", tdef.key, err)
			continue
		}

		transition := &models.TransitionNodeDetail{
			NodeID:      node.ID,
			Title:       tdef.title,
			Mode:        tdef.mode,
			Description: tdef.desc,
			CreatedAt:   time.Now(),
		}
		if err := db.WithContext(ctx).Create(transition).Error; err != nil {
			log.Printf("Warning: Could not create transition detail %s: %v", tdef.key, err)
			continue
		}

		nodeLibrary[tdef.key] = node.ID
	}

	// STEP 2: Link nodes to plans using nodeLibrary
	// Structure: Plan -> [Attraction, Transition, Attraction, Transition, ...] (no consecutive attractions)
	planNodeSequences := map[string][]struct {
		nodeKey      string // reference to nodeLibrary
		planNote     string // plan-specific routing/tips
		priceCents   int
		durationMins int
	}{
		"Summer European Adventure": {
			{nodeKey: "paris-eiffel", planNote: "Go early morning to beat crowds. Don't miss the top floor.", priceCents: 2500, durationMins: 120},
			{nodeKey: "flight-intl", planNote: "Book budget airlines. Arrive 2 hours early.", priceCents: 8000, durationMins: 180},
			{nodeKey: "rome-colosseum", planNote: "Buy skip-the-line tickets. Bring water and sunscreen.", priceCents: 1800, durationMins: 90},
			{nodeKey: "walk-short", planNote: "Explore historic streets to next attraction.", priceCents: 1, durationMins: 45},
			{nodeKey: "rome-vatican", planNote: "Dress code: shoulders and knees covered. Allow extra time.", priceCents: 3000, durationMins: 180},
			{nodeKey: "train-overnight", planNote: "Book sleeper cabin for comfort. Great views through France.", priceCents: 12000, durationMins: 2880},
			{nodeKey: "barcelona-sagrada", planNote: "Book tickets in advance. Audio guide worth it.", priceCents: 2600, durationMins: 120},
		},
		"Tokyo Food Tour": {
			{nodeKey: "tokyo-tsukiji", planNote: "Go early morning for freshest options. Try tamagoyaki.", priceCents: 3000, durationMins: 180},
			{nodeKey: "subway-tokyo", planNote: "Get Suica/Pasmo card for easy transit.", priceCents: 300, durationMins: 15},
			{nodeKey: "tokyo-ramen", planNote: "Pick any shop, all excellent. Tonkotsu is a must-try.", priceCents: 1500, durationMins: 90},
			{nodeKey: "walk-short", planNote: "Scenic walk through Meiji Shrine forest.", priceCents: 1, durationMins: 45},
			{nodeKey: "tokyo-meiji", planNote: "Beautiful peaceful shrine in forested area.", priceCents: 1, durationMins: 60},
			{nodeKey: "taxi-shibuya", planNote: "Quick ride to Shibuya for nightlife.", priceCents: 1500, durationMins: 15},
			{nodeKey: "tokyo-shibuya", planNote: "Experience the world's busiest crossing.", priceCents: 1, durationMins: 30},
		},
		"New Zealand Road Trip": {
			{nodeKey: "nz-milford", planNote: "Book cruise in advance. Bring rain jacket.", priceCents: 5000, durationMins: 240},
			{nodeKey: "drive-scenic", planNote: "Most scenic drive in world. Stop often for photos.", priceCents: 4000, durationMins: 240},
			{nodeKey: "nz-queenstown", planNote: "Book bungy or jet boat. Worth the adrenaline!", priceCents: 15000, durationMins: 180},
			{nodeKey: "walk-short", planNote: "Scenic walk by the lake with amazing views.", priceCents: 1, durationMins: 90},
			{nodeKey: "nz-glow-worms", planNote: "Mind-blowing glow ceiling. No photography.", priceCents: 3500, durationMins: 120},
			{nodeKey: "flight-regional", planNote: "Collect rental car at airport.", priceCents: 6000, durationMins: 120},
			{nodeKey: "nz-hobbiton", planNote: "Spend full day. Tours are comprehensive.", priceCents: 5000, durationMins: 180},
		},
		"Budget Southeast Asia Backpacking": {
			{nodeKey: "bangkok-palace", planNote: "Dress respectfully - shoulders and knees covered.", priceCents: 500, durationMins: 120},
			{nodeKey: "flight-regional", planNote: "Book with AirAsia - very affordable.", priceCents: 2500, durationMins: 120},
			{nodeKey: "phuket-beach", planNote: "Stay in backpacker hostels. Amazing food at cheap prices.", priceCents: 1, durationMins: 240},
			{nodeKey: "bus-transit", planNote: "Budget option. Bring neck pillow and copy of passport.", priceCents: 3000, durationMins: 960},
			{nodeKey: "hanoi-oldquarter", planNote: "Try pho and egg coffee. Negotiations required at markets.", priceCents: 1500, durationMins: 180},
			{nodeKey: "ferry", planNote: "Budget tour boat. Splurge is worth it for experience.", priceCents: 5000, durationMins: 1440},
			{nodeKey: "angkor-wat", planNote: "Hire guide for history lessons. Watch sunrise from temple.", priceCents: 3700, durationMins: 480},
		},
		"Caribbean Island Hopping": {
			{nodeKey: "cayman-beach", planNote: "Best beach in Caribbean. Rent snorkel gear.", priceCents: 1, durationMins: 480},
			{nodeKey: "ferry", planNote: "Fast and comfortable. Bring passport.", priceCents: 5000, durationMins: 120},
			{nodeKey: "jamaica-falls", planNote: "Climb terraces carefully. Water is fresh and cool.", priceCents: 2500, durationMins: 120},
			{nodeKey: "flight-regional", planNote: "Small plane but quick hop.", priceCents: 7500, durationMins: 90},
			{nodeKey: "turks-grace", planNote: "Even more pristine than Seven Mile.", priceCents: 1, durationMins: 240},
			{nodeKey: "bus-transit", planNote: "See dolphins and tropical fish.", priceCents: 4500, durationMins: 300},
		},
		"Machu Picchu and Peruvian Highlands": {
			{nodeKey: "cusco-city", planNote: "Essential for avoiding altitude sickness. Visit local markets.", priceCents: 2000, durationMins: 240},
			{nodeKey: "train-express", planNote: "Book PeruRail. Splurge on first class for views.", priceCents: 9000, durationMins: 300},
			{nodeKey: "machu-picchu", planNote: "Sunrise hike worth the 4am wake-up. Hire guide for history.", priceCents: 15000, durationMins: 480},
			{nodeKey: "walk-short", planNote: "Steep but well-maintained path. Knee-friendly route.", priceCents: 1, durationMins: 180},
			{nodeKey: "cusco-market", planNote: "Perfect for souvenirs. Local vendors are friendly.", priceCents: 3000, durationMins: 120},
		},
		"Iceland Ring Road Adventure": {
			{nodeKey: "iceland-gullfoss", planNote: "Bring rain jacket. Epic scale of water flow.", priceCents: 1, durationMins: 60},
			{nodeKey: "drive-scenic", planNote: "Stop at waterfalls and black sand beaches.", priceCents: 6000, durationMins: 300},
			{nodeKey: "iceland-black-sand", planNote: "Powerful waves. Stay back from water.", priceCents: 1, durationMins: 90},
			{nodeKey: "drive-highway", planNote: "Rough roads in winter. Essential 4WD.", priceCents: 4500, durationMins: 240},
			{nodeKey: "iceland-geysir", planNote: "Strokkur geyser erupts every 5-10 minutes.", priceCents: 1500, durationMins: 120},
			{nodeKey: "walk-short", planNote: "Scenic walk through geothermal areas.", priceCents: 1, durationMins: 45},
			{nodeKey: "iceland-blue-lagoon", planNote: "Book in advance. Silica mud mask is included.", priceCents: 6500, durationMins: 180},
		},
		"Norway Fjords Road Trip": {
			{nodeKey: "drive-scenic", planNote: "One of world's most scenic drives. Stop often for photos.", priceCents: 5000, durationMins: 300},
			{nodeKey: "ferry", planNote: "Experience Norwegian fjords by water. Spectacular views.", priceCents: 3000, durationMins: 180},
			{nodeKey: "walk-short", planNote: "Hike through mountain passes with alpine views.", priceCents: 1, durationMins: 120},
			{nodeKey: "train-express", planNote: "Scenic train through mountains. Book first class.", priceCents: 8000, durationMins: 240},
			{nodeKey: "walk-short", planNote: "Walk to waterfall viewpoint.", priceCents: 1, durationMins: 60},
		},
		"Morocco Cultural Immersion": {
			{nodeKey: "flight-intl", planNote: "Direct flight to Casablanca. Easy connection to Marrakech.", priceCents: 9000, durationMins: 180},
			{nodeKey: "taxi", planNote: "Traditional taxi through medina streets.", priceCents: 500, durationMins: 30},
			{nodeKey: "hanoi-oldquarter", planNote: "Wander souks. Haggle for local crafts and spices.", priceCents: 2000, durationMins: 240},
			{nodeKey: "walk-short", planNote: "Walk through blue-painted alleys of chefchaouen.", priceCents: 1, durationMins: 90},
			{nodeKey: "bus-transit", planNote: "Local bus to Sahara desert camp.", priceCents: 2000, durationMins: 480},
			{nodeKey: "bike", planNote: "Camel trek in Sahara. Sunset is magical.", priceCents: 8000, durationMins: 360},
		},
		"India Taj Mahal and Beyond": {
			{nodeKey: "flight-intl", planNote: "Arrive in Delhi. Book accommodation near Old Delhi.", priceCents: 12000, durationMins: 240},
			{nodeKey: "train-express", planNote: "Overnight express to Agra. Book AC tier for comfort.", priceCents: 4000, durationMins: 600},
			{nodeKey: "machu-picchu", planNote: "Sunrise visit mandatory. Hire guide for history and photography tips.", priceCents: 3000, durationMins: 180},
			{nodeKey: "drive-scenic", planNote: "Road trip through Rajasthan. Stop at villages.", priceCents: 3000, durationMins: 240},
			{nodeKey: "walk-short", planNote: "Walk through markets and ancient forts.", priceCents: 1, durationMins: 90},
			{nodeKey: "taxi", planNote: "Taxi to temple for evening rituals and prayer ceremony.", priceCents: 500, durationMins: 30},
		},
	}

	// STEP 2: Link nodes to plans using nodeLibrary
	// Create PlanNode records linking generic nodes to plans with plan-specific context
	planNodeMap := make(map[string][]string) // planTitle -> []nodeIDs for later use

	for planTitle, nodeSequence := range planNodeSequences {
		planID, exists := planMap[planTitle]
		if !exists {
			log.Printf("Warning: Plan %s not found, skipping node linking", planTitle)
			continue
		}

		seqPos := 1
		var lastNodeType string // Track previous node type for validation

		for _, seq := range nodeSequence {
			nodeID, exists := nodeLibrary[seq.nodeKey]
			if !exists {
				log.Printf("Warning: Node reference %s not found for plan %s", seq.nodeKey, planTitle)
				continue
			}

			// Get node type for validation
			node := &models.Node{}
			if err := db.WithContext(ctx).Model(&models.Node{}).Where("id = ?", nodeID).First(node).Error; err != nil {
				log.Printf("Warning: Could not fetch node type for %s", seq.nodeKey)
				continue
			}

			// Validate: No two consecutive attractions (spec requirement)
			if lastNodeType == "attraction" && node.Type == "attraction" {
				log.Printf("Warning: Skipping consecutive attraction %s in plan %s at position %d", seq.nodeKey, planTitle, seqPos)
				continue
			}
			lastNodeType = node.Type

			// Create PlanNode record with plan-specific context
			planNode := &models.PlanNode{
				ID:                  uuid.New().String(),
				PlanID:              planID,
				NodeID:              nodeID,
				SequencePosition:    seqPos,
				Description:         &seq.planNote, // Plan-specific routing tips
				EstimatedPriceCents: &seq.priceCents,
				DurationMinutes:     &seq.durationMins,
				CreatedAt:           time.Now(),
			}

			if err := db.WithContext(ctx).Create(planNode).Error; err != nil {
				log.Printf("Warning: Could not link node %s to plan %s: %v", seq.nodeKey, planTitle, err)
				continue
			}

			planNodeMap[planTitle] = append(planNodeMap[planTitle], nodeID)
			seqPos++
			log.Printf("✓ Linked node [%d] %s to plan %s", seqPos-1, seq.nodeKey, planTitle)
		}
	}

	// Step 3: Create ratings for plans
	ratingData := []struct {
		planTitle string
		userID    string
		stars     int
	}{
		{
			planTitle: "Summer European Adventure",
			userID:    userMap["bob"],
			stars:     5,
		},
		{
			planTitle: "Tokyo Food Tour",
			userID:    userMap["alice"],
			stars:     4,
		},
		{
			planTitle: "New Zealand Road Trip",
			userID:    userMap["bob"],
			stars:     5,
		},
		{
			planTitle: "Budget Southeast Asia Backpacking",
			userID:    userMap["charlie"],
			stars:     4,
		},
		{
			planTitle: "Machu Picchu and Peruvian Highlands",
			userID:    userMap["alice"],
			stars:     5,
		},
	}

	for _, rd := range ratingData {
		planID, exists := planMap[rd.planTitle]
		if !exists {
			continue
		}

		// Create rating
		rating := &models.Rating{
			ID:        uuid.New().String(),
			PlanID:    planID,
			UserID:    rd.userID,
			Stars:     rd.stars,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.WithContext(ctx).Create(rating).Error; err != nil {
			log.Printf("Warning: Could not seed rating for plan %s: %v", rd.planTitle, err)
			continue
		}

		log.Printf("✓ Seeded rating: %d stars for plan %s", rd.stars, rd.planTitle)
	}

	// Create comments on plans
	commentData := []struct {
		planTitle string
		authorID  string
		text      string
	}{
		{
			planTitle: "Summer European Adventure",
			authorID:  userMap["bob"],
			text:      "What's the best time of year to do this trip?",
		},
		{
			planTitle: "Tokyo Food Tour",
			authorID:  userMap["charlie"],
			text:      "I'm vegetarian - are there good options?",
		},
		{
			planTitle: "New Zealand Road Trip",
			authorID:  userMap["alice"],
			text:      "How long did this take? Planning for summer vacation.",
		},
		{
			planTitle: "Budget Southeast Asia Backpacking",
			authorID:  userMap["bob"],
			text:      "This is perfect for my backpacking trip!",
		},
		{
			planTitle: "Machu Picchu and Peruvian Highlands",
			authorID:  userMap["charlie"],
			text:      "Do I need to book guides in advance?",
		},
	}

	for _, cd := range commentData {
		planID, exists := planMap[cd.planTitle]
		if !exists {
			continue
		}

		now := time.Now()
		comment := &models.Comment{
			ID:        uuid.New().String(),
			PlanID:    planID,
			AuthorID:  cd.authorID,
			Text:      cd.text,
			CreatedAt: now,
			UpdatedAt: &now,
		}

		if err := db.WithContext(ctx).Create(comment).Error; err != nil {
			log.Printf("Warning: Could not seed comment for plan %s: %v", cd.planTitle, err)
			continue
		}

		// Increment comment count on the plan
		if err := db.WithContext(ctx).Model(&models.TravelPlan{}).Where("id = ?", planID).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
			log.Printf("Warning: Could not update plan comment count for plan %s: %v", cd.planTitle, err)
		}

		log.Printf("✓ Seeded comment on plan %s", cd.planTitle)
	}

	log.Println("Database seeding completed successfully!")
	return nil
}
