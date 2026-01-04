module.exports = {
  preset: 'jest-preset-angular',
  setupFilesAfterEnv: ['<rootDir>/setup-jest.ts'],
  testPathIgnorePatterns: [
    '<rootDir>/node_modules/',
    '<rootDir>/dist/',
    '<rootDir>/src/open-api/'
  ],
  // Performance optimizations
  maxWorkers: '50%',
  cache: true,
  cacheDirectory: '<rootDir>/.jest-cache',
  // Transform configuration with isolatedModules for faster compilation
  transform: {
    '^.+\\.(ts|js|mjs|html|svg)$': [
      'jest-preset-angular',
      {
        tsconfig: '<rootDir>/tsconfig.spec.json',
        stringifyContentPathRegex: '\\.(html|svg)$',
      }
    ]
  },
  collectCoverageFrom: [
    'src/**/*.ts',
    '!src/**/*.d.ts',
    '!src/main.ts',
    '!src/open-api/**',
    '!src/environments/**'
  ],
  coverageDirectory: 'coverage/receipt-wrangler-desktop',
  coverageReporters: ['html', 'lcov', 'cobertura', 'text-summary'],
  coverageThreshold: {
    global: {
      statements: 1,
      branches: 1,
      functions: 1,
      lines: 1
    }
  },
  moduleNameMapper: {
    '^src/(.*)$': '<rootDir>/src/$1'
  },
  transformIgnorePatterns: [
    'node_modules/(?!.*\\.mjs$|@angular|@ngxs|rxjs|ngx-mask|@ng-bootstrap|ngx-bootstrap|pretty-print-json)'
  ],
  // Additional performance settings
  testEnvironmentOptions: {
    customExportConditions: ['node', 'node-addons']
  },
  workerIdleMemoryLimit: '512MB'
};
