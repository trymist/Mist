
export interface BaseEntity {
  id: number;
  createdAt?: string;
  updatedAt?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  hasNext: boolean;
  hasPrev: boolean;
}

export interface SelectOption {
  label: string;
  value: string;
}

export interface FormField {
  name: string;
  label: string;
  type: 'text' | 'email' | 'password' | 'textarea' | 'select' | 'tags';
  required?: boolean;
  options?: SelectOption[];
  defaultValue?: string | string[] | number;
}

export interface AppError {
  message: string;
  code?: string;
  field?: string;
}
