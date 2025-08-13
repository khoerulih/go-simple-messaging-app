## Menggunakan image dasar Golang versi 1.24 berbasis Alpine Linux (ringan dan cepat)
FROM golang:1.24-alpine

## Menetapkan direktori kerja di dalam container ke /app
WORKDIR /app

## Menyalin file dependency Go (go.mod & go.sum) ke dalam container
COPY go.mod go.sum ./

## Mengunduh semua dependency yang dibutuhkan sesuai go.mod & go.sum
RUN go mod tidy

## Menyalin seluruh source code aplikasi ke dalam container
COPY . .

## Build aplikasi Go menjadi binary bernama simple-messaging-app
RUN go build -o simple-messaging-app

## Memberikan permission eksekusi pada binary hasil build
RUN chmod +x simple-messaging-app

## Membuka port 4000 dan 8080 pada container (untuk akses aplikasi)
EXPOSE 4000
EXPOSE 8080

## Menjalankan aplikasi saat container dijalankan
CMD [ "./simple-messaging-app" ]