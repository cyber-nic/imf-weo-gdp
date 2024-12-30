import type { D3Node, DataRow } from './types';

/**
 * Converts flat data into a D3-compatible tree structure.
 * @param data - The flat data array from the CSV.
 * @param year - The selected year to use for values.
 * @returns A tree structure with a root node.
 */
export function parseDataToTree(data: DataRow[], year: string): D3Node {
	// Root node for the tree
	const root: D3Node = {
		name: year,
		children: []
	};

	// Regex to match a 3-letter country code at the start of SERIES_CODE
	const countryCodeRegex = /^[A-Z]{3}\./;

	// Populate children with country nodes
	root.children = data
		.filter((row) => countryCodeRegex.test(row['SERIES_CODE'])) // Include only rows with valid country codes
		.map((row) => {
			const countryName = row['COUNTRY.Name'].split(',')[0]; // Extract country name from full name
			const gdpValue = parseFloat(row[year]) || 0; // Default to 0 if GDP is missing

			return {
				name: countryName,
				value: gdpValue
			};
		})
		.filter((node) => node.value > 0); // Filter out nodes with no GDP

	return root;
}
