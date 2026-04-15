-- Migration: Seed system nodes (attractions and transitions)
-- Provides a pre-populated database with sample attractions and transitions
-- All system nodes are pre-approved (is_approved=true) for immediate use

-- First, ensure system user exists (for created_by foreign key)
INSERT OR IGNORE INTO users (id, email, username, password_hash, role, is_active, created_at)
VALUES (
    'system-user-001',
    'system@travellink.local',
    'system',
    '$2a$12$dummy.hash.for.system.user.that.cannot.be.used.for.login',
    'traveller',
    true,
    CURRENT_TIMESTAMP
);

-- Hotel Nodes (8 total)
INSERT OR IGNORE INTO nodes (id, type, created_by, is_approved, created_at) VALUES
('attraction-hotel-001', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-hotel-002', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-hotel-003', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-hotel-004', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP);

INSERT OR IGNORE INTO attraction_node_details (node_id, name, category, location, description, contact_info, hours_of_operation) VALUES
('attraction-hotel-001', 'Grand Plaza Hotel', 'hotel', '100 Central Plaza, Downtown', 'Luxury 5-star hotel with world-class amenities. Full-service spa, fine dining, and conference facilities. 24-hour concierge service.', '+1-555-0201', '24/7'),
('attraction-hotel-002', 'Comfort Suites Express', 'hotel', '200 Highway 101, Business Park', 'Affordable 3-star hotel perfect for business travelers. Free breakfast, fitness center, and WiFi included.', '+1-555-0202', '24/7'),
('attraction-hotel-003', 'Waterfront Residences', 'hotel', '300 Harbor Drive, Waterfront District', 'Boutique 4-star hotel with ocean views. Modern rooms with smart technology and eco-friendly amenities.', '+1-555-0203', '24/7'),
('attraction-hotel-004', 'Mountain Lodge & Resort', 'hotel', '400 Peak Road, Mountain View', 'Scenic resort set in natural surroundings. Hiking trails, outdoor activities, and rustic elegance. WiFi available.', '+1-555-0204', '24/7');

-- Museum Nodes (6 total)
INSERT OR IGNORE INTO nodes (id, type, created_by, is_approved, created_at) VALUES
('attraction-museum-001', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-museum-002', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-museum-003', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-museum-004', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-museum-005', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-museum-006', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP);

INSERT OR IGNORE INTO attraction_node_details (node_id, name, category, location, description, contact_info, hours_of_operation) VALUES
('attraction-museum-001', 'National History Museum', 'museum', '1000 Knowledge Avenue, Cultural Center', 'Extensive collection of fossils, minerals, and natural history artifacts. Interactive exhibits for all ages. IMAX theater on-site.', '+1-555-0301', '09:00-18:00'),
('attraction-museum-002', 'Modern Art Gallery', 'museum', '1100 Creativity Street, Arts Quarter', 'Contemporary art exhibitions featuring local and international artists. Regular workshops and artist talks. Sculpture garden included.', '+1-555-0302', '10:00-19:00'),
('attraction-museum-003', 'Science & Innovation Center', 'museum', '1200 Tech Boulevard, Innovation Hub', 'Interactive science exhibits including planetarium, virtual reality experiences, and hands-on experiments.', '+1-555-0303', '09:30-17:30'),
('attraction-museum-004', 'Historical Heritage Museum', 'museum', '1300 Old Town Road, Historic District', 'Preserves local history with artifacts, documents, and multimedia presentations. Regular guided tours available.', '+1-555-0304', '10:00-16:00'),
('attraction-museum-005', 'Maritime Museum', 'museum', '1400 Dock Street, Harbor Area', 'Historic ships, nautical artifacts, and maritime history. Walk-through submarine replica and interactive navigation exhibits.', '+1-555-0305', '09:00-17:00'),
('attraction-museum-006', 'Art of the Ancient World', 'museum', '1500 Classical Lane, Downtown', 'Egyptian, Greek, and Roman artifacts. Authentic sculptures, pottery, and jewelry spanning thousands of years.', '+1-555-0306', '10:00-18:00');

-- Shopping Nodes (5 total)
INSERT OR IGNORE INTO nodes (id, type, created_by, is_approved, created_at) VALUES
('attraction-shopping-001', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-shopping-002', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-shopping-003', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-shopping-004', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-shopping-005', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP);

INSERT OR IGNORE INTO attraction_node_details (node_id, name, category, location, description, contact_info, hours_of_operation) VALUES
('attraction-shopping-001', 'Central Shopping Mall', 'shopping', '2000 Commerce Drive, Commercial Zone', 'Multi-level shopping center with 200+ stores. Features major international and luxury brands. Food court and cinemas included.', '+1-555-0401', '09:00-21:00'),
('attraction-shopping-002', 'Artisan Market Square', 'shopping', '2100 Craft Lane, Arts District', 'Curated local vendors selling handmade crafts, art, and unique merchandise. Weekly live performances and food trucks.', '+1-555-0402', '08:00-18:00'),
('attraction-shopping-003', 'Vintage & Antique District', 'shopping', '2200 Memory Lane, Historic Quarter', 'Multiple vintage shops and antique dealers. Collectibles, vintage clothing, and rare items. Treasure hunting experience.', '+1-555-0403', '10:00-17:00'),
('attraction-shopping-004', 'Fashion Boulevard', 'shopping', '2300 Style Street, Fashion District', 'Luxury boutiques and designer flagship stores. High-end fashion, watches, and accessories. Personal shopping services available.', '+1-555-0404', '10:00-20:00'),
('attraction-shopping-005', 'Farmers Market Fresh', 'shopping', '2400 Garden Road, Riverside', 'Weekly farmers market with local produce, honey, artisan breads, and handmade goods. Peak hours: Saturday & Sunday.', '+1-555-0405', '06:00-14:00');

-- Parks and Recreation (5 total)
INSERT OR IGNORE INTO nodes (id, type, created_by, is_approved, created_at) VALUES
('attraction-park-001', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-park-002', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-park-003', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-park-004', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-park-005', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP);

INSERT OR IGNORE INTO attraction_node_details (node_id, name, category, location, description, contact_info, hours_of_operation) VALUES
('attraction-park-001', 'Central City Park', 'park', '3000 Park Avenue, City Center', 'Urban oasis with 200 acres. Walking trails, duck pond, botanical gardens, and picnic areas. Free admission.', '+1-555-0501', '06:00-22:00'),
('attraction-park-002', 'Riverside Nature Preserve', 'park', '3100 Nature Trail, Countryside', 'Protected wetlands and wildlife habitat. Birdwatching opportunities, nature photography and guided ecological tours.', '+1-555-0502', '08:00-17:00'),
('attraction-park-003', 'Adventure Sports Park', 'park', '3200 Active Road, Recreation Zone', 'Ziplines, rock climbing, obstacle courses, and adventure activities. Professional instructors and safety equipment provided.', '+1-555-0503', '09:00-18:00'),
('attraction-park-004', 'Beach & Boardwalk', 'park', '3300 Coastal Path, Beach District', 'Sandy beach with boardwalk attractions. Swimming, volleyball courts, beach bars. Popular sunset viewing spot.', '+1-555-0504', '06:00-20:00'),
('attraction-park-005', 'Children''s Playground & Zoo', 'park', '3400 Family Lane, Educational Zone', 'Interactive children''s park with play structures. Adjacent small zoo with local wildlife and petting areas.', '+1-555-0505', '09:00-17:00');

-- Tourist Attractions (6 total)
INSERT OR IGNORE INTO nodes (id, type, created_by, is_approved, created_at) VALUES
('attraction-landmark-001', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-landmark-002', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-landmark-003', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-landmark-004', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-landmark-005', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-landmark-006', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP);

INSERT OR IGNORE INTO attraction_node_details (node_id, name, category, location, description, contact_info, hours_of_operation) VALUES
('attraction-landmark-001', 'Iconic Clock Tower', 'tourist_attraction', '4000 Historic Square, Downtown', 'Restored 1890s clock tower offering panoramic city views. Historical architecture tour and gift shop available.', '+1-555-0601', '10:00-18:00'),
('attraction-landmark-002', 'Botanical Gardens', 'tourist_attraction', '4100 Green Drive, Nature Quarter', '30 acres of curated gardens featuring plants from around the world. Conservatory and cafe with garden views.', '+1-555-0602', '08:00-19:00'),
('attraction-landmark-003', 'Observation Deck Tower', 'tourist_attraction', '4200 Sky Road, Uptown', 'Visit the highest point in the city. 360-degree views, restaurant, gift shop, and interactive exhibits.', '+1-555-0603', '09:00-22:00'),
('attraction-landmark-004', 'Historic Bridge Tour', 'tourist_attraction', '4300 Riverside, Waterfront', 'Walk across iconic suspension bridge with scenic overlooks. Guided tours explain engineering and local history.', '+1-555-0604', '09:00-17:00'),
('attraction-landmark-005', 'Street Art Mural District', 'tourist_attraction', '4400 Creative Way, Arts District', 'Walking tour of world-famous street art murals and graffiti galleries. Photography tours available.', '+1-555-0605', '10:00-18:00'),
('attraction-landmark-006', 'Entertainment Theater District', 'tourist_attraction', '4500 Broadway, Theater Quarter', 'Historic theaters hosting live performances, musicals, and concerts. Guided backstage tours available.', '+1-555-0606', '10:00-22:00');

-- Entertainment Venues (4 total)
INSERT OR IGNORE INTO nodes (id, type, created_by, is_approved, created_at) VALUES
('attraction-entertainment-001', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-entertainment-002', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-entertainment-003', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP),
('attraction-entertainment-004', 'attraction', 'system-user-001', true, CURRENT_TIMESTAMP);

INSERT OR IGNORE INTO attraction_node_details (node_id, name, category, location, description, contact_info, hours_of_operation) VALUES
('attraction-entertainment-001', 'Concert & Music Venue', 'entertainment', '5000 Sound Boulevard, Arts District', 'Intimate venue hosting local and international musicians. Acoustically excellent for all music genres.', '+1-555-0701', '19:00-23:00'),
('attraction-entertainment-002', 'Comedy Club & Bar', 'entertainment', '5100 Laugh Lane, Downtown', 'Nightly comedy shows with food and full bar service. Mix of established and up-and-coming comedians.', '+1-555-0702', '19:00-22:00'),
('attraction-entertainment-003', 'Bowling & Arcade Alley', 'entertainment', '5200 Fun Street, Recreation Zone', 'Modern bowling lanes, arcade games, billiards, and casual dining. Perfect for group outings.', '+1-555-0703', '12:00-23:00'),
('attraction-entertainment-004', 'Dance Club & Lounge', 'entertainment', '5300 Groove Road, Party District', 'Upscale nightclub with DJ, dance floor, and cocktail lounge. Weekend special events and themed parties.', '+1-555-0704', '21:00-02:00');

-- ============================================================================
-- TRANSITION NODES (Transportation modes between attractions)
-- ============================================================================

-- Walking Transitions (5 total)
INSERT OR IGNORE INTO nodes (id, type, created_by, is_approved, created_at) VALUES
('transition-walking-001', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-walking-002', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-walking-003', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-walking-004', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-walking-005', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP);

INSERT OR IGNORE INTO transition_node_details (node_id, title, mode, description, hours_of_operation) VALUES
('transition-walking-001', 'Walking', 'walking', 'Short walk through downtown streets', NULL),
('transition-walking-002', 'Walking', 'walking', 'Scenic walk along waterfront promenade', NULL),
('transition-walking-003', 'Walking', 'walking', 'Moderate walk through business district', NULL),
('transition-walking-004', 'Walking', 'walking', 'Leisurely walk through historic neighborhood', NULL),
('transition-walking-005', 'Walking', 'walking', 'Quick walk to nearby attractions in same district', NULL);

-- Driving Transitions (5 total)
INSERT OR IGNORE INTO nodes (id, type, created_by, is_approved, created_at) VALUES
('transition-driving-001', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-driving-002', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-driving-003', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-driving-004', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-driving-005', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP);

INSERT OR IGNORE INTO transition_node_details (node_id, title, mode, description, hours_of_operation) VALUES
('transition-driving-001', 'Personal Car', 'car', 'Drive across downtown via major highways', NULL),
('transition-driving-002', 'Personal Car', 'car', 'Scenic drive to mountain resort (with traffic)', NULL),
('transition-driving-003', 'Personal Car', 'car', 'Interstate drive to neighboring city', NULL),
('transition-driving-004', 'Personal Car', 'car', 'Local surface streets through residential area', NULL),
('transition-driving-005', 'Personal Car', 'car', 'Express drive to airport or transit hub', NULL);

-- Public Transit (Bus, Train) - 6 total
INSERT OR IGNORE INTO nodes (id, type, created_by, is_approved, created_at) VALUES
('transition-transit-001', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-transit-002', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-transit-003', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-transit-004', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-transit-005', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-transit-006', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP);

INSERT OR IGNORE INTO transition_node_details (node_id, title, mode, description, hours_of_operation) VALUES
('transition-transit-001', 'Line 5 Downtown Bus', 'bus', 'City bus with multiple stops through downtown', 'Daily 6:00-23:00'),
('transition-transit-002', 'Express Line 22', 'bus', 'Express bus to regional shopping center', 'Daily 7:00-22:00'),
('transition-transit-003', 'Commuter Rail RT100', 'train', 'Commuter rail to neighboring districts', 'Mon-Fri 5:30-23:00, Sat-Sun 7:00-22:00'),
('transition-transit-004', 'Metro Express Line M1', 'train', 'Express metro service underground', 'Daily 5:00-1:00 (next day)'),
('transition-transit-005', 'Airport Shuttle A1', 'bus', 'Airport shuttle service', 'Daily 5:00-23:30'),
('transition-transit-006', 'Light Rail Loop LRL', 'train', 'Light rail loop through entertainment district', 'Daily 6:00-midnight');

-- Taxi/Rideshare - 3 total
INSERT OR IGNORE INTO nodes (id, type, created_by, is_approved, created_at) VALUES
('transition-taxi-001', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-taxi-002', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-taxi-003', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP);

INSERT OR IGNORE INTO transition_node_details (node_id, title, mode, description, hours_of_operation) VALUES
('transition-taxi-001', 'Taxi/Rideshare', 'taxi', 'Taxi ride across city center with light traffic', 'Daily 24/7'),
('transition-taxi-002', 'Rideshare Service', 'taxi', 'Rideshare service to remote attractions', 'Daily 24/7'),
('transition-taxi-003', 'Premium Car Service', 'taxi', 'Premium car service to upscale location', 'Daily 24/7');

-- Biking - 4 total
INSERT OR IGNORE INTO nodes (id, type, created_by, is_approved, created_at) VALUES
('transition-bike-001', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-bike-002', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-bike-003', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-bike-004', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP);

INSERT OR IGNORE INTO transition_node_details (node_id, title, mode, description, hours_of_operation) VALUES
('transition-bike-001', 'Bike Path - Park', 'bike', 'Bike ride on dedicated bike path through park', NULL),
('transition-bike-002', 'Mountain Bike Trail', 'bike', 'Mountain bike trail on hilly terrain', NULL),
('transition-bike-003', 'Urban Bike Route', 'bike', 'Urban bike commute with protected lanes', NULL),
('transition-bike-004', 'Riverside Bike Path', 'bike', 'Scenic riverside bike path', NULL);

-- Miscellaneous Transitions - 3 total
INSERT OR IGNORE INTO nodes (id, type, created_by, is_approved, created_at) VALUES
('transition-other-001', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-other-002', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP),
('transition-other-003', 'transition', 'system-user-001', true, CURRENT_TIMESTAMP);

INSERT OR IGNORE INTO transition_node_details (node_id, title, mode, description, hours_of_operation) VALUES
('transition-other-001', 'River Ferry', 'other', 'Ferry service across river or bay', 'Daily 7:00-19:00'),
('transition-other-002', 'Historic Cable Car', 'other', 'Cable car up steep hill with views', 'Daily 9:00-17:00'),
('transition-other-003', 'Hotel Shuttle', 'other', 'Free shuttle service between attractions', 'Daily 8:00-22:00');
