import { BaseLevel } from './base_level'; 


export interface BaseBenefit {
    id: string;
    name: string;
    type: string;
    frequency: string;
    discount: number;
    stock: number;
    isActive: boolean;
    level: BaseLevel;
    numTimes: number;
    createdAt: string;
    updatedAt: string;
    externalProductId: string;
    dependency: boolean;
    minAmount: number
}


