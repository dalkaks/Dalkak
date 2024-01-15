import '@emotion/react';

type ColorId = 'gray50' | 'gray100' | 'gray200' | 'gray400';

// type TypographyId =
//   | 'bigTitle'
//   | 'h1'
//   | 'h2'
//   | 'h3'
//   | 'content1'
//   | 'content2'
//   | 'content3'
//   | 'content4'
//   | 'bigTitleBold'
//   | 'h1Bold'
//   | 'h2Bold'
//   | 'h3Bold'
//   | 'content1Bold'
//   | 'content2Bold'
//   | 'content3Bold'
//   | 'content4Bold'
//   | 'buttonXS'
//   | 'buttonS'
//   | 'buttonM'
//   | 'menuButton';

declare module '@emotion/react' {
	export interface Theme {
		color: {
			[key in ColorId]: string;
		};
		// typography: {
		//   [key in TypographyId]: {
		//     fontSize: string;
		//     lineHeight?: string;
		//     letterSpacing: string;
		//     fontWeight: number;
		//   };
		// };
	}
}
