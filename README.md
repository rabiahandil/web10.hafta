# GoLearn - Uzaktan Eğitim Platformu API

GoLearn, öğretmenlerin kurs ve içerik oluşturabildiği, öğrencilerin ise bu kurslara katılıp ilerlemelerini takip edebildiği modern bir backend API projesidir.

## Özellikler
- **Kimlik Doğrulama:** JWT tabanlı güvenli giriş ve kayıt sistemi.
- **RBAC (Role-Based Access Control):** Öğretmen ve Öğrenci rolleri için farklı yetkilendirme seviyeleri.
- **Kurs & İçerik Yönetimi:** Kurs, ders ve quiz oluşturma (sadece öğretmenler).
- **İlerleme Takibi:** Ders tamamlama ve kurs bazlı ilerleme yüzdesi hesaplama.
- **Quiz Sistemi:** Çoktan seçmeli quizler, otomatik puanlama ve sonuç kaydı.
- **Canlı Sınıf (WebSocket):** Her kurs için özel odalarda gerçek zamanlı mesajlaşma.
- **Genişletilmiş Filtreleme:** Kurslarda kategoriye göre filtreleme ve sayfalama.
- **Güvenlik:** IP bazlı hız sınırlama (Rate Limiting).
- **Dokümantasyon:** Swagger (OpenAPI 3.0) entegrasyonu.

## Teknoloji Yığını
- **Go (Golang)** - Gin Framework
- **GORM** - SQLite (Veritabanı)
- **JWT** - Kimlik doğrulama
- **Gorilla WebSocket** - Gerçek zamanlı mesajlaşma
- **Docker** - Konteynerizasyon

## Çalıştırma

### Docker ile (Önerilen)
```bash
docker-compose up --build
```
Uygulama `http://localhost:8090` adresinde çalışacaktır.

### Yerel Olarak
1. Bağımlılıkları yükleyin: `go mod download`
2. Swagger dokümantasyonunu oluşturun (isteğe bağlı): `swag init`
3. Uygulamayı çalıştırın: `go run main.go`

## API Dokümantasyonu
Swagger arayüzüne şu adresten erişebilirsiniz:
`http://localhost:8090/swagger/index.html`

## WebSocket Kullanımı
Bağlantı URL'si: `ws://localhost:8090/ws/classroom/{course_id}`
Not: Bağlantı sırasında `Authorization: Bearer {token}` header'ı veya token doğrulama gereklidir.

## Örnek Endpoint'ler

### Auth
- `POST /api/auth/register` - Kayıt ol
- `POST /api/auth/login` - Giriş yap (Token döner)

### Courses (Auth Gerekli)
- `GET /api/courses` - Tüm kursları listele (Filter: page, limit, category)
- `POST /api/courses` - Yeni kurs oluştur (Sadece Teacher)
- `PUT /api/courses/:id` - Kursu güncelle (Sadece Kurs Sahibi)

### Progress
- `POST /api/lessons/:id/complete` - Dersi tamamlandı olarak işaretle
- `GET /api/my/progress` - Genel ilerlemeni gör
