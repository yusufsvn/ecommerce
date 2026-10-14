import { Box, Container, Typography, Stack, Button, Avatar, Paper } from '@mui/material';
function About() {
  return (
    <Box sx={{ bgcolor: '#f5f5f5', minHeight: '100vh', pb: 8 }}>
      
      {/* 1. Hero Bölümü (Başlık) */}
      <Box sx={{ bgcolor: 'primary.main', color: 'white', py: 8, mb: 6 }}>
        <Container maxWidth="lg">
          <Stack spacing={2} alignItems="center" textAlign="center">
            <Typography variant="h2" component="h1" fontWeight="bold">
              Biz Kimiz?
            </Typography>
            <Typography variant="h6" sx={{ opacity: 0.9, maxWidth: '800px' }}>
              Teknoloji ve tasarımı birleştirerek dijital dünyada iz bırakıyoruz. 
              Müşterilerimize değer katmak için buradayız.
            </Typography>
          </Stack>
        </Container>
      </Box>

      <Container maxWidth="lg">
        
        {/* 2. Hikayemiz ve Görsel Bölümü */}
        <Stack 
          direction={{ xs: 'column', md: 'row' }} 
          spacing={6} 
          alignItems="center" 
          sx={{ mb: 8 }}
        >
          {/* Sol Taraf: Metin */}
          <Box sx={{ flex: 1 }}>
            <Stack spacing={3}>
              <Typography variant="h4" color="primary" fontWeight="bold">
                Hikayemiz & Misyonumuz
              </Typography>
              <Typography variant="body1" color="text.secondary" lineHeight={1.8}>
                2020 yılında küçük bir ofiste başlayan yolculuğumuz, bugün global ölçekte hizmet veren 
                bir teknoloji stüdyosuna dönüştü. Amacımız, karmaşık problemleri basit ve estetik 
                arayüzlerle çözmektir.
              </Typography>
              <Typography variant="body1" color="text.secondary" lineHeight={1.8}>
                Sadece kod yazmıyoruz; markanızın dijital kimliğini inşa ediyoruz. Kullanıcı deneyimini 
                merkeze alan yaklaşımımızla, sürdürülebilir ve ölçeklenebilir projeler geliştiriyoruz.
              </Typography>
              <Box>
                <Button variant="contained" size="large">
                  İletişime Geç
                </Button>
              </Box>
            </Stack>
          </Box>

          {/* Sağ Taraf: Görsel */}
          <Box sx={{ flex: 1, width: '100%' }}>
            <Box 
              component="img"
              src="https://images.unsplash.com/photo-1522071820081-009f0129c71c?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80"
              alt="Ofis ortamında çalışan ekip"
              sx={{
                width: '100%',
                height: 'auto',
                borderRadius: 4,
                boxShadow: 3,
                objectFit: 'cover'
              }}
            />
          </Box>
        </Stack>

        {/* 3. İstatistikler (Info Cards) */}
        <Stack 
          direction={{ xs: 'column', md: 'row' }} 
          spacing={4} 
          justifyContent="center" 
          sx={{ mb: 8 }}
        >
          {[
            { title: 'Tamamlanan Proje', count: '150+' },
            { title: 'Mutlu Müşteri', count: '80+' },
            { title: 'Deneyim Yılı', count: '5+' },
          ].map((item, index) => (
            <Paper 
              key={index} 
              elevation={2} 
              sx={{ 
                p: 4, 
                flex: 1, 
                textAlign: 'center', 
                borderRadius: 2 
              }}
            >
              <Typography variant="h3" color="primary" fontWeight="bold">
                {item.count}
              </Typography>
              <Typography variant="subtitle1" color="text.secondary">
                {item.title}
              </Typography>
            </Paper>
          ))}
        </Stack>

        {/* 4. Ekibimiz Bölümü */}
        <Box sx={{ mb: 4 }}>
          <Typography variant="h4" textAlign="center" fontWeight="bold" sx={{ mb: 6 }}>
            Lider Ekibimiz
          </Typography>
          
          <Stack 
            direction={{ xs: 'column', md: 'row' }} 
            spacing={4} 
            justifyContent="center"
          >
            {/* Kart 1 */}
            <Stack alignItems="center" spacing={2} sx={{ flex: 1 }}>
              <Avatar 
                src="https://images.unsplash.com/photo-1560250097-0b93528c311a?auto=format&fit=crop&w=400&q=80" 
                sx={{ width: 150, height: 150, mb: 1 }}
              />
              <Typography variant="h6" fontWeight="bold">Ahmet Yılmaz</Typography>
              <Typography variant="body2" color="text.secondary">Kurucu & CEO</Typography>
            </Stack>

            {/* Kart 2 */}
            <Stack alignItems="center" spacing={2} sx={{ flex: 1 }}>
              <Avatar 
                src="https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?auto=format&fit=crop&w=400&q=80" 
                sx={{ width: 150, height: 150, mb: 1 }}
              />
              <Typography variant="h6" fontWeight="bold">Zeynep Kaya</Typography>
              <Typography variant="body2" color="text.secondary">Baş Tasarımcı</Typography>
            </Stack>

            {/* Kart 3 */}
            <Stack alignItems="center" spacing={2} sx={{ flex: 1 }}>
              <Avatar 
                src="https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?auto=format&fit=crop&w=400&q=80" 
                sx={{ width: 150, height: 150, mb: 1 }}
              />
              <Typography variant="h6" fontWeight="bold">Mehmet Demir</Typography>
              <Typography variant="body2" color="text.secondary">Yazılım Lideri</Typography>
            </Stack>
          </Stack>
        </Box>

      </Container>
    </Box>
  )
}

export default About