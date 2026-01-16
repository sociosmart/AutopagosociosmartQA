export interface BasePagination<T> {
    nextCursor?: string;
    prevCursor?: string;
    total: number;
    items: T[];
}