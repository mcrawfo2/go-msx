const { jsWithBabel: tsJestPreset } = require('ts-jest/presets');

module.exports = {
	preset: 'ts-jest',
	globals: {
		'ts-jest': {
			babelConfig: 'babel.config.js',
			diagnostics: false
		}
	},
	reporters: [
		'default',
    'jest-junit'
	],
	projects: [
		{
			displayName: 'test',
			testMatch: [
        '<rootDir>/src/**/*.spec.ts',
			],
			moduleFileExtensions: ['ts', 'js'],
			moduleNameMapper: {
				// Replace .scss files with empty modules during testing
				'\\.scss$': '<rootDir>/src/spec-helpers/empty-module.js',

        '@msx/http': '<rootDir>/src/spec-helpers/empty-module.js',
				'^spec-helpers/(.*)': '<rootDir>/src/spec-helpers/$1'
			},
			setupFilesAfterEnv: [
				'<rootDir>/jest.init.js'
			],
			transform: {
				...tsJestPreset.transform,
				'\\.html$': '<rootDir>/src/spec-helpers/transformers/html.js'
			},
			transformIgnorePatterns: [
				'node_modules/(?!vms-)'
			]
		}
	],
  coverageDirectory: "test"
};
