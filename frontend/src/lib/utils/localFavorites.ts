import { writable } from 'svelte/store';

const LOCAL_FAVORITES_KEY = 'marchive_local_favorites';

export interface LocalFavorite {
	bookHash: string;
	addedAt: number;
}

function createLocalFavoritesStore() {
	const { subscribe, set, update } = writable<LocalFavorite[]>([]);

	// Initialize from localStorage
	if (typeof window !== 'undefined') {
		try {
			const stored = localStorage.getItem(LOCAL_FAVORITES_KEY);
			if (stored) {
				set(JSON.parse(stored));
			}
		} catch {
			set([]);
		}
	}

	return {
		subscribe,
		add: (bookHash: string) => {
			update(favorites => {
				if (!favorites.find(f => f.bookHash === bookHash)) {
					const newFavorites = [...favorites, { bookHash, addedAt: Date.now() }];
					if (typeof window !== 'undefined') {
						localStorage.setItem(LOCAL_FAVORITES_KEY, JSON.stringify(newFavorites));
					}
					return newFavorites;
				}
				return favorites;
			});
		},
		remove: (bookHash: string) => {
			update(favorites => {
				const newFavorites = favorites.filter(f => f.bookHash !== bookHash);
				if (typeof window !== 'undefined') {
					localStorage.setItem(LOCAL_FAVORITES_KEY, JSON.stringify(newFavorites));
				}
				return newFavorites;
			});
		},
		isFavorited: (bookHash: string) => {
			let result = false;
			update(favorites => {
				result = favorites.some(f => f.bookHash === bookHash);
				return favorites;
			});
			return result;
		},
		getCount: () => {
			let count = 0;
			update(favorites => {
				count = favorites.length;
				return favorites;
			});
			return count;
		}
	};
}

export const localFavorites = createLocalFavoritesStore();
