import { BaseBenefitGenerated } from './base_benefit_generated';
import { BaseCustomer } from './base_customer';

export interface BaseBenefitTicket  {
    id: string;
    customer: BaseCustomer;
    startDate: string;
    endDate: string;
    benefitGenerated: BaseBenefitGenerated;
    redeemed: boolean;
    createdAt: string;
    updatedAt: string;
}


