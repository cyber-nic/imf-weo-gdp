export interface D3Node {
	name: string;
	value?: number;
	children?: D3Node[];
	x0?: number;
	x1?: number;
	y0?: number;
	y1?: number;
}

// Define the type for a single row in the CSV data
export interface DataRow {
	[key: string]: string;
	'COUNTRY.Name': string;
}
