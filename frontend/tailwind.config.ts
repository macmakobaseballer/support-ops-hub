import type { Config } from 'tailwindcss'

export default {
  content: [
    './components/**/*.{js,vue,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './plugins/**/*.{js,ts}',
    './app.vue',
  ],
  theme: {
    extend: {
      colors: {
        primary:   '#005461',
        secondary: '#018790',
        accent:    '#00B7B5',
        'page-bg': '#F4F4F4',
        sidebar: {
          bg:      '#005461',
          hover:   '#018790',
          active:  '#00B7B5',
          text:    'rgba(255,255,255,0.6)',
          section: 'rgba(255,255,255,0.35)',
        },
        status: {
          'new-bg':          '#dbeafe',
          'new-text':        '#1d4ed8',
          'in-progress-bg':  '#f5ede4',
          'in-progress-text':'#7c4d33',
          'waiting-bg':      '#ede9fe',
          'waiting-text':    '#7c3aed',
          'done-bg':         '#dcfce7',
          'done-text':       '#166534',
        },
        priority: {
          'high-bg':     '#fee2e2',
          'high-text':   '#dc2626',
          'medium-bg':   '#fef3c7',
          'medium-text': '#b45309',
          'low-bg':      '#f0fdf4',
          'low-text':    '#166534',
        },
      },
      spacing: {
        sidebar: '240px',
      },
      fontFamily: {
        sans: ['-apple-system', 'BlinkMacSystemFont', '"Segoe UI"', 'sans-serif'],
      },
      boxShadow: {
        card:       '0 1px 3px rgba(0,0,0,0.07)',
        'card-hover': '0 4px 12px rgba(0,0,0,0.12)',
        login:      '0 20px 60px rgba(0,0,0,0.3)',
      },
      borderRadius: {
        card: '10px',
      },
    },
  },
  plugins: [],
} satisfies Config
