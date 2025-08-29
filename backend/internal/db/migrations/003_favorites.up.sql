-- Add favorites table for authenticated users
CREATE TABLE favorites (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    book_hash TEXT NOT NULL,
    created_at BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_hash) REFERENCES savedbooks(hash) ON DELETE CASCADE,
    UNIQUE(user_id, book_hash)
);

-- Indexes for favorites
CREATE INDEX idx_favorites_user_id ON favorites(user_id);
CREATE INDEX idx_favorites_book_hash ON favorites(book_hash);
CREATE INDEX idx_favorites_created_at ON favorites(created_at);
