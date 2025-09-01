import js from '@eslint/js'
import pluginVue from 'eslint-plugin-vue'
import * as parserVue from 'vue-eslint-parser'
import configPrettier from '@vue/eslint-config-prettier'
import globals from 'globals' // You'll need to install this

export default [
  {
    name: 'app/files-to-lint',
    files: ['**/*.{js,mjs,jsx,vue}'],
  },

  {
    name: 'app/files-to-ignore',
    ignores: ['**/dist/**', '**/dist-ssr/**', '**/coverage/**'],
  },

  js.configs.recommended,
  ...pluginVue.configs['flat/recommended'],

  {
    name: 'app/vue-rules',
    languageOptions: {
      parser: parserVue,
      parserOptions: {
        ecmaVersion: 'latest',
        sourceType: 'module',
      },
      globals: {
        ...globals.browser, // Adds all browser globals including localStorage, fetch, etc.
        ...globals.node, // Adds process for your config files
      },
    },
    rules: {
      // Essential rules for error prevention
      'vue/no-unused-components': 'error',
      'vue/no-unused-vars': 'error',
      'vue/no-use-v-if-with-v-for': 'error',
      'vue/require-v-for-key': 'error',

      // Code quality rules
      'vue/no-undef-components': [
        'error',
        {
          ignorePatterns: [
            'router-view',
            'router-link',
            'keep-alive',
            'transition',
            'transition-group',
          ],
        },
      ],
      'vue/no-potential-component-option-typo': 'warn',
      'vue/prefer-true-attribute-shorthand': 'warn',

      // Style/readability rules
      'vue/component-name-in-template-casing': ['warn', 'PascalCase'],
      'vue/prop-name-casing': ['warn', 'camelCase'],
      'vue/v-on-event-hyphenation': ['warn', 'always'],

      // Relaxed rules (turn off overly strict ones)
      'vue/multi-word-component-names': 'off',
      'vue/require-default-prop': 'off',

      // Standard JavaScript rules
      'no-unused-vars': 'error',
      'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
      'no-debugger': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    },
  },

  configPrettier,
]
