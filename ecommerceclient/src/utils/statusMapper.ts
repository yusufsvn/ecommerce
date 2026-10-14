export const statusOptions = ['Bekliyor', 'Hazırlanıyor', 'Kargoda', 'Teslim Edildi', 'İptal Edildi'];

export const formatStatusToTurkish = (status: string) => {
  switch (status?.toLowerCase()) {
    case 'delivered':
      return 'Teslim Edildi';
    case 'shipped':
      return 'Kargoda';
    case 'preparing':
      return 'Hazırlanıyor';
    case 'cancelled':
      return 'İptal Edildi';
    case 'pending':
    default:
      return 'Bekliyor';
  }
};

export const mapTurkishToEnglishStatus = (status: string) => {
  if (!status) return '';
  switch (status) {
    case 'Teslim Edildi':
      return 'delivered';
    case 'Kargoda':
      return 'shipped';
    case 'Hazırlanıyor':
      return 'preparing';
    case 'İptal Edildi':
      return 'cancelled';
    case 'Bekliyor':
      return 'pending';
    default:
      return status;
  }
};

export const getStatusColor = (status: string) => {
  switch (status) {
    case 'Teslim Edildi':
      return 'success';
    case 'Kargoda':
      return 'info';
    case 'Hazırlanıyor':
      return 'warning';
    case 'İptal Edildi':
      return 'error';
    case 'Bekliyor':
    default:
      return 'default';
  }
};

// Ürün Yayın Durumu Mappers (API için UPPERCASE ENGLISH)
export const publishStatusOptions = ['Yayınla', 'Yayından Kaldır'];

export const mapTurkishToEnglishPublishState = (state: string): string => {
  if (!state) return 'UNPUBLISHED';
  const normalized = state.trim().toUpperCase();
  if (normalized === 'YAYINLA' || normalized === 'YAYINLANDI' || normalized === 'PUBLISHED') {
    return 'PUBLISHED';
  }
  if (normalized === 'YAYINDAN KALDIR' || normalized === 'YAYINDAN KALDIRILDI' || normalized === 'UNPUBLISHED') {
    return 'UNPUBLISHED';
  }
  return normalized;
};

export const mapEnglishToTurkishPublishState = (state: string): string => {
  if (!state) return 'Yayından Kaldır';
  const normalized = state.trim().toUpperCase();
  if (normalized === 'PUBLISHED') {
    return 'Yayınla';
  }
  if (normalized === 'UNPUBLISHED') {
    return 'Yayından Kaldır';
  }
  return state;
};

// Ürün Tipi Mappers (API için UPPERCASE ENGLISH)
export const itemTypeOptions = ['Fiziksel', 'Sanal'];

export const mapTurkishToEnglishItemType = (type: string): string => {
  if (!type) return 'PHYSICAL';
  const normalized = type.trim().toUpperCase();
  if (normalized === 'FIZIKSEL' || normalized === 'FİZİKSEL' || normalized === 'PHYSICAL') {
    return 'PHYSICAL';
  }
  if (normalized === 'SANAL' || normalized === 'VIRTUAL') {
    return 'VIRTUAL';
  }
  return normalized;
};

export const mapEnglishToTurkishItemType = (type: string): string => {
  if (!type) return 'Fiziksel';
  const normalized = type.trim().toUpperCase();
  if (normalized === 'PHYSICAL') {
    return 'Fiziksel';
  }
  if (normalized === 'VIRTUAL') {
    return 'Sanal';
  }
  return type;
};

