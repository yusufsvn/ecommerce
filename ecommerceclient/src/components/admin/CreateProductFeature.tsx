import React, { useState } from 'react';
import {
    Box,
    TextField,
    Button,
    List,
    ListItem,
    ListItemText,
    IconButton,
    Typography,
    Stack
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';

export interface ProductFeature {
    key: string;
    value: string;
}

export interface ProductFeatureInputProps {
    features: ProductFeature[]; // Üst bileşenden gelen liste
    onChange: (features: ProductFeature[]) => void; // Üst bileşene veri gönderecek fonksiyon
}
function ProductFeatureInput(props: ProductFeatureInputProps) {

    const [key, setKey] = useState<string>('')
    const [value, setValue] = useState<string>('')


    // klavyeden yapılan her girişi alır state e kaydeder
    const onChangeKey = (event: React.ChangeEvent<HTMLInputElement>) => {
        setKey(event.target.value)
    }
    const onChangeValue = (event: React.ChangeEvent<HTMLInputElement>) => {
        setValue(event.target.value)
    }

    const handleAdd = () => {
        if (!key.trim() || !value.trim()) {
            return;
        }

        // save input
        const newFeature: ProductFeature = {
            key: key.trim(),
            value: value.trim()
        }

        props.onChange([...props.features, newFeature])
        // clear form 
        setKey('')
        setValue('')
    }

    const handleRemove = (index: number) => {
        const temp = props.features.filter((_, i) => index !== i)
        props.onChange(temp)
    }


    return (
        <Box sx={{ maxWidth: 600 }}>
            <Typography variant="h6" gutterBottom>
                Product Features
            </Typography>

            <Stack direction="row" spacing={2} sx={{ mb: 2 }}>
                <TextField
                    label="Key (e.g. Color)"
                    variant="outlined"
                    value={key}
                    size="small"
                    fullWidth
                    onChange={onChangeKey}
                />
                <TextField
                    label="Value (e.g. Red)"
                    variant="outlined"
                    size="small"
                    fullWidth
                    value={value}
                    onChange={onChangeValue}
                />
                <Button variant="contained" onClick={handleAdd} disableElevation>
                    Add
                </Button>
            </Stack>

            <List>
                {props.features.map((feature, index) => (
                    <ListItem
                        key={index}
                        secondaryAction={
                            <IconButton edge="end" aria-label="delete" onClick={() => handleRemove(index)}>
                                <DeleteIcon color="error" />
                            </IconButton>
                        }
                        sx={{ bgcolor: 'background.paper', mb: 1, border: '1px solid #e0e0e0', borderRadius: 1 }}
                    >
                        <ListItemText
                            primary={feature.key}
                            secondary={feature.value}
                        />
                    </ListItem>
                ))}
            </List>
        </Box>
    );
};


export default ProductFeatureInput