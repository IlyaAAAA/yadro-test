# Запуск локально
go run cmd/app.go "filename"

# Запуск docker контейнера
docker build -t yadro_test .
docker run -it yadro_test test_file_1.txt

или .\start.sh