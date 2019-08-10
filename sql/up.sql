CREATE TABLE IF NOT EXISTS orders (
  sequence_id INT AUTO_INCREMENT,
  id VARCHAR(50),
  status VARCHAR(255) NOT NULL,
  distance FLOAT(10, 2) DEFAULT 0.0,
  location_info TEXT NOT NULL,
  PRIMARY KEY (sequence_id),
  UNIQUE (id)
);