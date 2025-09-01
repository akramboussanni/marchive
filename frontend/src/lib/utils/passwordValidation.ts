export interface PasswordValidationResult {
	isValid: boolean;
	errors: string[];
	requirementsText: string;
}

export function validatePassword(password: string): PasswordValidationResult {
	const errors: string[] = [];
	const requirements: string[] = [];

	// Check minimum length
	if (password.length < 8) {
		errors.push('At least 8 characters');
	} else {
		requirements.push('At least 8 characters');
	}

	// Check for lowercase letter
	if (!/[a-z]/.test(password)) {
		errors.push('One lowercase letter');
	} else {
		requirements.push('One lowercase letter');
	}

	// Check for uppercase letter
	if (!/[A-Z]/.test(password)) {
		errors.push('One uppercase letter');
	} else {
		requirements.push('One uppercase letter');
	}

	// Check for number
	if (!/\d/.test(password)) {
		errors.push('One number');
	} else {
		requirements.push('One number');
	}

	// Check for special character
	if (!/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password)) {
		errors.push('One special character');
	} else {
		requirements.push('One special character');
	}

	return {
		isValid: errors.length === 0,
		errors,
		requirementsText: requirements.join(', ')
	};
}

export function getPasswordStrength(password: string): 'weak' | 'medium' | 'strong' {
	if (!password) return 'weak';

	let score = 0;

	// Length contribution
	if (password.length >= 8) score += 1;
	if (password.length >= 12) score += 1;
	if (password.length >= 16) score += 1;

	// Character variety contribution
	if (/[a-z]/.test(password)) score += 1;
	if (/[A-Z]/.test(password)) score += 1;
	if (/\d/.test(password)) score += 1;
	if (/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password)) score += 1;

	// Bonus for mixed case and numbers
	if (/[a-z]/.test(password) && /[A-Z]/.test(password)) score += 1;
	if (/\d/.test(password) && /[a-zA-Z]/.test(password)) score += 1;

	// Determine strength based on score
	if (score >= 6) return 'strong';
	if (score >= 4) return 'medium';
	return 'weak';
}
