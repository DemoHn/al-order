CREATE TABLE IF NOT EXISTS orders (
  id INT AUTO_INCREMENT,
  status VARCHAR(255) NOT NULL,
  distance FLOAT(10, 2),
  origin_lat FLOAT(10, 2) NOT NULL,
  origin_lng FLOAT(10, 2) NOT NULL,
  dest_lat FLOAT(10, 2) NOT NULL,
  dest_lng FLOAT(10, 2) NOT NULL,
  PRIMARY KEY (id)
);