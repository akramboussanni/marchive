<script lang="ts">
	import { validatePassword, getPasswordStrength, type PasswordValidationResult } from '$lib/utils/passwordValidation';

	export let id: string;
	export let label: string;
	export let placeholder: string = '';
	export let value: string = '';
	export let disabled: boolean = false;
	export let required: boolean = false;
	export let showRequirements: boolean = false;
	export let showStrength: boolean = true;
	export let showMatchIndicator: boolean = false;
	export let confirmPassword: string = '';

	// Password validation state
	let passwordValidation: PasswordValidationResult | null = null;
	let passwordStrength: 'weak' | 'medium' | 'strong' = 'weak';
	let showPasswordRequirements = false;

	// Reactive password validation
	$: if (value) {
		passwordValidation = validatePassword(value);
		passwordStrength = getPasswordStrength(value);
	} else {
		passwordValidation = null;
		passwordStrength = 'weak';
	}

	// Check if passwords match (when confirmPassword is provided)
	$: passwordsMatch = showMatchIndicator && confirmPassword && value === confirmPassword;

	function getStrengthColor(strength: 'weak' | 'medium' | 'strong'): string {
		switch (strength) {
			case 'weak': return 'text-red-400';
			case 'medium': return 'text-yellow-400';
			case 'strong': return 'text-green-400';
			default: return 'text-gray-400';
		}
	}

	function getStrengthBg(strength: 'weak' | 'medium' | 'strong'): string {
		switch (strength) {
			case 'weak': return 'bg-red-500';
			case 'medium': return 'bg-yellow-500';
			case 'strong': return 'bg-green-500';
			default: return 'bg-gray-500';
		}
	}

	// Expose validation result to parent
	export { passwordValidation, passwordStrength, passwordsMatch };
</script>

<div>
	<label for={id} class="block text-sm font-medium text-gray-300 mb-2">
		{label}
	</label>
	<input
		{id}
		name={id}
		type="password"
		autocomplete="new-password"
		{required}
		bind:value
		{placeholder}
		{disabled}
		class="w-full px-4 py-3 border border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 bg-dark-700 text-gray-100 placeholder-gray-400"
	/>
	
	<!-- Password Requirements -->
	{#if showRequirements}
		<div class="mt-2">
			<button
				type="button"
				on:click={() => showPasswordRequirements = !showPasswordRequirements}
				class="text-xs text-gray-400 hover:text-gray-300 transition-colors"
			>
				{showPasswordRequirements ? 'Hide' : 'Show'} password requirements
			</button>
			
			{#if showPasswordRequirements}
				<div class="mt-2 p-3 bg-dark-700 rounded-lg border border-gray-600">
					<p class="text-xs text-gray-300 mb-2">Password must have:</p>
					<ul class="text-xs text-gray-400 space-y-1">
						{#if passwordValidation}
							{#each passwordValidation.errors as error}
								<li class="flex items-center space-x-2">
									<span class="text-red-400">✗</span>
									<span>{error}</span>
								</li>
							{/each}
							{#each passwordValidation.requirementsText.split(', ') as req}
								{#if !passwordValidation.errors.some(e => e.includes(req))}
									<li class="flex items-center space-x-2">
										<span class="text-green-400">✓</span>
										<span>{req}</span>
									</li>
								{/if}
							{/each}
						{:else}
							<li class="flex items-center space-x-2">
								<span class="text-gray-500">•</span>
								<span>At least 8 characters</span>
							</li>
							<li class="flex items-center space-x-2">
								<span class="text-gray-500">•</span>
								<span>One lowercase letter</span>
							</li>
							<li class="flex items-center space-x-2">
								<span class="text-gray-500">•</span>
								<span>One uppercase letter</span>
							</li>
							<li class="flex items-center space-x-2">
								<span class="text-gray-500">•</span>
								<span>One number</span>
							</li>
						{/if}
					</ul>
				</div>
			{/if}
		</div>
	{/if}

	<!-- Password Strength Indicator -->
	{#if showStrength && value}
		<div class="mt-2">
			<div class="flex items-center space-x-2">
				<span class="text-xs text-gray-400">Strength:</span>
				<span class="text-xs font-medium {getStrengthColor(passwordStrength)}">
					{passwordStrength.charAt(0).toUpperCase() + passwordStrength.slice(1)}
				</span>
			</div>
			<div class="mt-1 flex space-x-1">
				<div class="h-1 flex-1 rounded-full bg-gray-700">
					<div class="h-1 rounded-full {getStrengthBg(passwordStrength)} transition-all duration-300" 
						 style="width: {passwordStrength === 'weak' ? '33%' : passwordStrength === 'medium' ? '66%' : '100%'}"></div>
				</div>
			</div>
		</div>
	{/if}

	<!-- Password Match Indicator -->
	{#if showMatchIndicator && confirmPassword}
		<div class="mt-2">
			{#if passwordsMatch}
				<span class="text-xs text-green-400 flex items-center space-x-1">
					<span>✓</span>
					<span>Passwords match</span>
				</span>
			{:else}
				<span class="text-xs text-red-400 flex items-center space-x-1">
					<span>✗</span>
					<span>Passwords do not match</span>
				</span>
			{/if}
		</div>
	{/if}
</div>
