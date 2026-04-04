import js from '@eslint/js'; // Импорт базовых правил ESLint (для чистого JS)
import reactPlugin from 'eslint-plugin-react'; // Импорт плагина для React
import reactHooksPlugin from 'eslint-plugin-react-hooks'; // Импорт плагина для React Hooks
import globals from 'globals'; // Импорт глобальных переменных (для браузера и ES2021)

export default [
  js.configs.recommended, // Использование базовых правил ESLint
  {
    files: ['**/*.{js,jsx}'],
    plugins: {
      react: reactPlugin, // Использование плагина для React
      'react-hooks': reactHooksPlugin, // Использование плагина для React Hooks
    },
    languageOptions: {
      ecmaVersion: 'latest', // Использование последней версии ECMAScript
      sourceType: 'module', // Использование модульной системы
      globals: {
        ...globals.browser, // Использование глобальных переменных для браузера
        ...globals.es2021, // Использование глобальных переменных для ES2021
      },
      parserOptions: {
        ecmaFeatures: { jsx: true }, // Использование JSX
      },
    },
    settings: {
      react: { version: 'detect' }, // автоопределение версии React из package.json
    },
    rules: {
      'react/react-in-jsx-scope': 'off', // В React 17+ не нужно импортировать React (jsx атоватически преобразуется)
      'react/prop-types': 'off', // Отключение проверки PropTypes (вместо этого используем TypeScript)
      'react-hooks/rules-of-hooks': 'error', // Правила для React Hooks
      'react-hooks/exhaustive-deps': 'warn', // Проверка зависимостей в React Hooks
      'no-unused-vars': ['warn', { argsIgnorePattern: '^_' }], // Предупреждение о неиспользуемых переменных
      'no-console': ['warn', { allow: ['warn', 'error'] }], // Предупреждение о console.log (разрешены warn и error)
    },
  },
];
