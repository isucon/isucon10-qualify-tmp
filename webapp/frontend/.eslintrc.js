module.exports = {
  parser: '@typescript-eslint/parser',
  extends: [
    'standard',
    'standard-react',
    'standard-with-typescript',
    "plugin:react-hooks/recommended"
  ],
  plugins: [],
  rules: {
    'react/react-in-jsx-scope': 'off',
    'react/prop-types': 0,
    '@typescript-eslint/strict-boolean-expressions': 0,
    '@typescript-eslint/explicit-function-return-type': 0
  },
  globals: {
    React: 'writable'
  },
  parserOptions: {
    project: './tsconfig.json'
  }
}
