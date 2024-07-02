# Wait for mysql1 to be ready
until docker exec mysql_container mysql -u root -proot_password -e "SELECT 1"; do
  echo "Waiting for mysql database connection..."
  sleep 5
done

# Create table in mysql1
docker exec mysql_container mysql -u root -proot_password -e "
CREATE TABLE server.server (
    id INT AUTO_INCREMENT PRIMARY KEY,
    status VARCHAR(255) NOT NULL
);"