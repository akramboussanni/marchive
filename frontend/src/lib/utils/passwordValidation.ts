export interface PasswordRequirements {
	minLength: number;
	requireLower: boolean;
	requireUpper: boolean;
	requireDigit: boolean;
	requireSpecial: boolean;
}

export const defaultPasswordRequirements: PasswordRequirements = {
	minLength: 8,
	requireLower: true,
	requireUpper: true,
	requireDigit: true,
	requireSpecial: false
};

export interface PasswordValidationResult {
	isValid: boolean;
	errors: string[];
	requirementsText: string;
}

export function validatePassword(password: string, requirements: PasswordRequirements = defaultPasswordRequirements): PasswordValidationResult {
	const errors: string[] = [];
	
	if (password.length < requirements.minLength) {
		errors.push(`Password must be at least ${requirements.minLength} characters long`);
	}
	
	if (requirements.requireLower && !/[a-z]/.test(password)) {
		errors.push('Password must contain at least one lowercase letter');
	}
	
	if (requirements.requireUpper && !/[A-Z]/.test(password)) {
		errors.push('Password must contain at least one uppercase letter');
	}
	
	if (requirements.requireDigit && !/\d/.test(password)) {
		errors.push('Password must contain at least one number');
	}
	
	if (requirements.requireSpecial && !/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password)) {
		errors.push('Password must contain at least one special character');
	}
	
	const requirementsText = getPasswordRequirementsText(requirements);
	
	return {
		isValid: errors.length === 0,
		errors,
		requirementsText
	};
}

export function getPasswordRequirementsText(requirements: PasswordRequirements = defaultPasswordRequirements): string {
	const reqs: string[] = [];
	
	reqs.push(`At least ${requirements.minLength} characters`);
	
	if (requirements.requireLower) {
		reqs.push('One lowercase letter');
	}
	
	if (requirements.requireUpper) {
		reqs.push('One uppercase letter');
	}
	
	if (requirements.requireDigit) {
		reqs.push('One number');
	}
	
	if (requirements.requireSpecial) {
		reqs.push('One special character');
	}
	
	return reqs.join(', ');
}

export function getPasswordStrength(password: string): 'weak' | 'medium' | 'strong' {
	if (password.length < 6) return 'weak';
	
	let score = 0;
	if (password.length >= 8) score++;
	if (/[a-z]/.test(password)) score++;
	if (/[A-Z]/.test(password)) score++;
	if (/\d/.test(password)) score++;
	if (/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password)) score++;
	
	if (score <= 2) return 'weak';
	if (score <= 4) return 'medium';
	return 'strong';
}
