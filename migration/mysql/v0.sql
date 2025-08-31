-- Create users table
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create orders table
CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    price DECIMAL(10, 2) NOT NULL,
    status ENUM('pending', 'processing', 'shipped', 'delivered', 'cancelled') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Insert sample users
INSERT INTO users (username, email) VALUES
('john_doe', 'john@example.com'),
('jane_smith', 'jane@example.com'),
('mike_wilson', 'mike@example.com'),
('sarah_jones', 'sarah@example.com'),
('david_brown', 'david@example.com');

-- Insert sample orders
INSERT INTO orders (user_id, product_name, quantity, price, status) VALUES
(1, 'Laptop Computer', 1, 999.99, 'delivered'),
(1, 'Wireless Mouse', 2, 29.99, 'delivered'),
(2, 'Smartphone', 1, 699.99, 'shipped'),
(2, 'Phone Case', 1, 19.99, 'shipped'),
(3, 'Tablet', 1, 399.99, 'processing'),
(3, 'Screen Protector', 2, 9.99, 'processing'),
(4, 'Headphones', 1, 149.99, 'pending'),
(5, 'Keyboard', 1, 79.99, 'cancelled'),
(1, 'Monitor', 1, 299.99, 'pending'),
(2, 'USB Cable', 3, 12.99, 'delivered');