import { api } from "../Interceptors/httpinterceptor";
import type { GetCategories } from "../models/ProductResponse";

const AdminUrl = "https://ecommerceapi.yusufsevinc.work/api/admin";
const publicUrl = "https://ecommerceapi.yusufsevinc.work/api/public";

export const categoryService = {
    // Get all categories from public endpoint
    getCategories: async (): Promise<GetCategories[]> => {
        const response = await api.get<GetCategories[]>(publicUrl + "/getcategories");
        return response.data;
    },

    // Create a new category (parameters sent via query string)
    createCategory: async (category_name: string): Promise<GetCategories> => {
        const response = await api.post<GetCategories>(AdminUrl + "/createcategory", null, {
            params: {
                category_name: category_name
            },
            withCredentials: true
        });
        return response.data;
    },

    // Update an existing category (parameters sent via query string)
    updateCategory: async (category_id: string, category_name: string): Promise<boolean> => {
        const response = await api.put(AdminUrl + "/updatecategory", null, {
            params: {
                category_id: category_id,
                category_name: category_name
            },
            withCredentials: true
        });
        if (response.status === 200) {
            return true;
        }
        return false;
    },

    // Delete a category
    deleteCategory: async (category_id: string): Promise<boolean> => {
        const response = await api.delete(AdminUrl + "/deletecategory", {
            params: {
                category_id: category_id
            },
            withCredentials: true
        });
        if (response.status === 200) {
            return true;
        }
        return false;
    }
};
