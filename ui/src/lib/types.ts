export interface FeatureFlag {
	id: string;
	name: string;
	description: string;
	tenant: string;
	enabled: boolean;
	created_at: string;
	created_by: string;
	updated_at: string;
	updated_by: string;
}

export interface CreateFlagPayload {
	name: string;
	description: string;
	enabled: boolean;
}

export interface UpdateFlagPayload {
	description: string;
	enabled: boolean;
}