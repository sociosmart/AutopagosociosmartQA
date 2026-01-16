import { BaseBenefit } from './base_benefit';

export interface BaseBenefitGenerated extends BaseBenefit {
    startDate: string;
    endDate: string;
    benefit: BaseBenefit;
}


