import { useEffect, useState } from 'react';
import {
    Box,
    Paper,
    Typography,
    Button,
    Stack,
    Container,
    TextField,
    InputAdornment,
    IconButton,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    CircularProgress,
    Tooltip
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import SearchIcon from '@mui/icons-material/Search';
import AddIcon from '@mui/icons-material/Add';
import { categoryService } from '../../Services/categoryService';
import type { GetCategories } from '../../models/ProductResponse';
import { toast } from 'react-toastify';
import Spinner from '../Spinner';

export default function AdminCategories() {
    const [categories, setCategories] = useState<GetCategories[]>([]);
    const [loading, setLoading] = useState(true);
    const [searchQuery, setSearchQuery] = useState('');

    // Modal state for Add/Edit
    const [openDialog, setOpenDialog] = useState(false);
    const [dialogMode, setDialogMode] = useState<'create' | 'edit'>('create');
    const [selectedCategory, setSelectedCategory] = useState<GetCategories | null>(null);
    const [categoryNameInput, setCategoryNameInput] = useState('');
    const [submitting, setSubmitting] = useState(false);

    // Fetch categories
    const fetchCategories = async () => {
        try {
            setLoading(true);
            const data = await categoryService.getCategories();
            setCategories(data || []);
        } catch (error) {
            console.error("Kategoriler yüklenirken hata oluştu:", error);
            toast.error("Kategoriler yüklenirken bir hata oluştu.");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchCategories();
    }, []);

    // Handle Delete
    const handleDelete = async (id: string, name: string) => {
        if (window.confirm(`"${name}" kategorisini silmek istediğinize emin misiniz?`)) {
            try {
                const success = await categoryService.deleteCategory(id);
                if (success) {
                    setCategories((prev) => prev.filter((c) => c.category_id !== id));
                    toast.success("Kategori başarıyla silindi.");
                } else {
                    toast.error("Kategori silinirken bir hata oluştu.");
                }
            } catch (error) {
                console.error("Kategori silme hatası:", error);
                toast.error("Kategori silinirken bir hata oluştu.");
            }
        }
    };

    // Handle Open Dialog (Create)
    const handleOpenCreate = () => {
        setDialogMode('create');
        setSelectedCategory(null);
        setCategoryNameInput('');
        setOpenDialog(true);
    };

    // Handle Open Dialog (Edit)
    const handleOpenEdit = (category: GetCategories) => {
        setDialogMode('edit');
        setSelectedCategory(category);
        setCategoryNameInput(category.category_name);
        setOpenDialog(true);
    };

    // Handle Close Dialog
    const handleCloseDialog = () => {
        setOpenDialog(false);
        setSelectedCategory(null);
        setCategoryNameInput('');
    };

    // Handle Submit (Create or Update)
    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        const trimmedName = categoryNameInput.trim();
        if (!trimmedName) {
            toast.warning("Kategori adı boş olamaz!");
            return;
        }

        try {
            setSubmitting(true);
            if (dialogMode === 'create') {
                const newCat = await categoryService.createCategory(trimmedName);
                if (newCat) {
                    setCategories((prev) => [...prev, newCat]);
                    toast.success("Kategori başarıyla oluşturuldu.");
                    handleCloseDialog();
                } else {
                    toast.error("Kategori oluşturulamadı.");
                }
            } else if (dialogMode === 'edit' && selectedCategory) {
                const success = await categoryService.updateCategory(selectedCategory.category_id, trimmedName);
                if (success) {
                    setCategories((prev) =>
                        prev.map((c) =>
                            c.category_id === selectedCategory.category_id
                                ? { ...c, category_name: trimmedName }
                                : c
                        )
                    );
                    toast.success("Kategori başarıyla güncellendi.");
                    handleCloseDialog();
                } else {
                    toast.error("Kategori güncellenemedi.");
                }
            }
        } catch (error) {
            console.error("Kategori kaydetme hatası:", error);
            toast.error("Kategori kaydedilirken bir hata oluştu.");
        } finally {
            setSubmitting(false);
        }
    };

    // Filter categories based on search
    const filteredCategories = categories.filter((c) =>
        c.category_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        c.category_id.toLowerCase().includes(searchQuery.toLowerCase())
    );

    if (loading) {
        return <Spinner message="Kategoriler yükleniyor..." />;
    }

    return (
        <Container maxWidth="xl" sx={{ mt: 2, mb: 4 }}>
            {/* Header Section */}
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4, flexWrap: 'wrap', gap: 2 }}>
                <Box>
                    <Typography variant="h5" fontWeight="bold" color="text.primary">
                        Kategori Yönetimi
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Ürün kategorilerini ekleyebilir, güncelleyebilir veya silebilirsiniz.
                    </Typography>
                </Box>
                <Button
                    variant="contained"
                    color="primary"
                    startIcon={<AddIcon />}
                    onClick={handleOpenCreate}
                    sx={{ px: 3, py: 1, borderRadius: 2 }}
                >
                    Yeni Kategori Ekle
                </Button>
            </Box>

            {/* Filter Section */}
            <Paper sx={{ p: 2, mb: 3 }}>
                <TextField
                    fullWidth
                    size="small"
                    placeholder="Kategori Adı veya ID ile Ara..."
                    variant="outlined"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    InputProps={{
                        startAdornment: (
                            <InputAdornment position="start">
                                <SearchIcon color="action" />
                            </InputAdornment>
                        ),
                    }}
                />
            </Paper>

            {/* Table Section */}
            <TableContainer component={Paper} sx={{ overflow: 'hidden' }}>
                <Table>
                    <TableHead sx={{ bgcolor: 'grey.50' }}>
                        <TableRow>
                            <TableCell sx={{ fontWeight: 'bold', width: '40%' }}>Kategori ID</TableCell>
                            <TableCell sx={{ fontWeight: 'bold', width: '40%' }}>Kategori Adı</TableCell>
                            <TableCell align="right" sx={{ fontWeight: 'bold', width: '20%' }}>İşlemler</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {filteredCategories.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={3} align="center" sx={{ py: 4 }}>
                                    <Typography color="text.secondary">
                                        {searchQuery ? "Arama kriterine uygun kategori bulunamadı." : "Henüz kategori bulunmuyor."}
                                    </Typography>
                                </TableCell>
                            </TableRow>
                        ) : (
                            filteredCategories.map((category) => (
                                <TableRow
                                    key={category.category_id}
                                    sx={{
                                        '&:hover': { bgcolor: 'grey.200/50', transition: 'background-color 0.2s' }
                                    }}
                                >
                                    <TableCell sx={{ fontFamily: 'monospace', color: 'text.secondary' }}>
                                        {category.category_id}
                                    </TableCell>
                                    <TableCell sx={{ fontWeight: 'medium' }}>
                                        {category.category_name}
                                    </TableCell>
                                    <TableCell align="right">
                                        <Stack direction="row" spacing={1} justifyContent="flex-end">
                                            <Tooltip title="Düzenle">
                                                <IconButton
                                                    color="primary"
                                                    size="small"
                                                    onClick={() => handleOpenEdit(category)}
                                                    sx={{ border: '1px solid', borderColor: 'primary.light' }}
                                                >
                                                    <EditIcon fontSize="small" />
                                                </IconButton>
                                            </Tooltip>
                                            <Tooltip title="Sil">
                                                <IconButton
                                                    color="error"
                                                    size="small"
                                                    onClick={() => handleDelete(category.category_id, category.category_name)}
                                                    sx={{ border: '1px solid', borderColor: 'error.light' }}
                                                >
                                                    <DeleteIcon fontSize="small" />
                                                </IconButton>
                                            </Tooltip>
                                        </Stack>
                                    </TableCell>
                                </TableRow>
                            ))
                        )}
                    </TableBody>
                </Table>
            </TableContainer>

            {/* Create/Edit Dialog */}
            <Dialog open={openDialog} onClose={handleCloseDialog} maxWidth="sm" fullWidth>
                <form onSubmit={handleSubmit}>
                    <DialogTitle fontWeight="bold">
                        {dialogMode === 'create' ? "Yeni Kategori Ekle" : "Kategori Düzenle"}
                    </DialogTitle>
                    <DialogContent>
                        <Box sx={{ pt: 2 }}>
                            {dialogMode === 'edit' && selectedCategory && (
                                <TextField
                                    label="Kategori ID"
                                    value={selectedCategory.category_id}
                                    disabled
                                    fullWidth
                                    margin="normal"
                                    variant="outlined"
                                    sx={{ mb: 2 }}
                                />
                            )}
                            <TextField
                                label="Kategori Adı"
                                placeholder="Örn. Elektronik, Giyim vb."
                                value={categoryNameInput}
                                onChange={(e) => setCategoryNameInput(e.target.value)}
                                required
                                fullWidth
                                autoFocus
                                margin="normal"
                                variant="outlined"
                            />
                        </Box>
                    </DialogContent>
                    <DialogActions sx={{ px: 3, pb: 3 }}>
                        <Button onClick={handleCloseDialog} disabled={submitting}>
                            İptal
                        </Button>
                        <Button
                            type="submit"
                            variant="contained"
                            disabled={submitting}
                            startIcon={submitting ? <CircularProgress size={20} color="inherit" /> : null}
                        >
                            {dialogMode === 'create' ? "Oluştur" : "Güncelle"}
                        </Button>
                    </DialogActions>
                </form>
            </Dialog>
        </Container>
    );
}
