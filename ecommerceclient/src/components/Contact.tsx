import { Box, Container, Typography, Stack, TextField, Button, Paper } from '@mui/material';
import LocationOnIcon from '@mui/icons-material/LocationOn';
import PhoneIcon from '@mui/icons-material/Phone';
import EmailIcon from '@mui/icons-material/Email';
import SendIcon from '@mui/icons-material/Send';
function Contact() {
  return (
    <Box sx={{ bgcolor: '#f5f5f5', minHeight: '100vh', pb: 8 }}>
      
      {/* 1. Başlık Alanı */}
      <Box sx={{ bgcolor: 'primary.main', color: 'white', py: 6, mb: 6 }}>
        <Container maxWidth="lg">
          <Stack spacing={2} alignItems="center" textAlign="center">
            <Typography variant="h3" component="h1" fontWeight="bold">
              Bize Ulaşın
            </Typography>
            <Typography variant="h6" sx={{ opacity: 0.9, maxWidth: '600px' }}>
              Sorularınız mı var? Projeniz hakkında konuşmak mı istiyorsunuz? 
              Aşağıdaki formu doldurun veya doğrudan bizi arayın.
            </Typography>
          </Stack>
        </Container>
      </Box>

      <Container maxWidth="lg">
        
        {/* 2. Ana İçerik: Bilgiler ve Form */}
        <Stack 
          direction={{ xs: 'column', md: 'row' }} 
          spacing={4} 
          justifyContent="space-between"
          sx={{ mb: 8 }}
        >
          
          {/* SOL TARAF: İletişim Bilgileri Kartı */}
          <Paper elevation={3} sx={{ p: 4, flex: 1, borderRadius: 2 }}>
            <Stack spacing={4} height="100%" justifyContent="center">
              <Box>
                <Typography variant="h5" color="primary" fontWeight="bold" gutterBottom>
                  İletişim Bilgileri
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Mesai saatleri içinde (09:00 - 18:00) bize her zaman ulaşabilirsiniz.
                </Typography>
              </Box>

              <Stack spacing={3}>
                {/* Adres */}
                <Stack direction="row" spacing={2} alignItems="center">
                  <Box sx={{ bgcolor: 'primary.light', p: 1.5, borderRadius: '50%', color: 'white' }}>
                    <LocationOnIcon />
                  </Box>
                  <Box>
                    <Typography variant="subtitle1" fontWeight="bold">Adresimiz</Typography>
                    <Typography variant="body2" color="text.secondary">
                      Teknoloji Plaza, Kat: 3, No: 42<br />Nilüfer, Bursa
                    </Typography>
                  </Box>
                </Stack>

                {/* Telefon */}
                <Stack direction="row" spacing={2} alignItems="center">
                  <Box sx={{ bgcolor: 'primary.light', p: 1.5, borderRadius: '50%', color: 'white' }}>
                    <PhoneIcon />
                  </Box>
                  <Box>
                    <Typography variant="subtitle1" fontWeight="bold">Telefon</Typography>
                    <Typography variant="body2" color="text.secondary">
                      +90 (224) 123 45 67
                    </Typography>
                  </Box>
                </Stack>

                {/* E-posta */}
                <Stack direction="row" spacing={2} alignItems="center">
                  <Box sx={{ bgcolor: 'primary.light', p: 1.5, borderRadius: '50%', color: 'white' }}>
                    <EmailIcon />
                  </Box>
                  <Box>
                    <Typography variant="subtitle1" fontWeight="bold">E-posta</Typography>
                    <Typography variant="body2" color="text.secondary">
                      info@sirketiniz.com
                    </Typography>
                  </Box>
                </Stack>
              </Stack>
            </Stack>
          </Paper>

          {/* SAĞ TARAF: İletişim Formu */}
          <Paper elevation={3} sx={{ p: 4, flex: 1.5, borderRadius: 2 }}>
            <Stack spacing={3}>
              <Typography variant="h5" color="primary" fontWeight="bold">
                Mesaj Gönderin
              </Typography>
              
              <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
                <TextField label="Adınız" variant="outlined" fullWidth />
                <TextField label="Soyadınız" variant="outlined" fullWidth />
              </Stack>
              
              <TextField label="E-posta Adresi" variant="outlined" type="email" fullWidth />
              <TextField label="Konu" variant="outlined" fullWidth />
              
              <TextField 
                label="Mesajınız" 
                variant="outlined" 
                multiline 
                rows={4} 
                fullWidth 
              />
              
              <Button 
                variant="contained" 
                size="large" 
                endIcon={<SendIcon />}
                sx={{ alignSelf: 'flex-start', px: 4, py: 1.5 }}
              >
                Gönder
              </Button>
            </Stack>
          </Paper>

        </Stack>

        {/* 3. Harita Bölümü (Google Maps Embed) */}
        <Box sx={{ width: '100%', height: '400px', borderRadius: 4, overflow: 'hidden', boxShadow: 3 }}>
          <iframe 
            src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d48152.56345846286!2d28.941574!3d40.220839!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x14ca1582d8d45695%3A0x61a00555c96109a2!2sBursa!5e0!3m2!1str!2str!4v1700000000000!5m2!1str!2str" 
            width="100%" 
            height="100%" 
            style={{ border: 0 }} 
            allowFullScreen={false} 
            loading="lazy" 
            referrerPolicy="no-referrer-when-downgrade"
            title="Ofis Konumu"
          ></iframe>
        </Box>

      </Container>
    </Box>
  )
}

export default Contact