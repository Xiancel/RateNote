CREATE TABLE items (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    comment TEXT,
    rating DOUBLE PRECISION NOT NULL CHECK (rating >= 1 AND rating <= 10),
    image_path TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);