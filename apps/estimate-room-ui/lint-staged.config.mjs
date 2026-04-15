const runUiChecks = () => [
  'npm run lint',
  'npm run typecheck',
  'npm test -- --run'
];

export default {
  '**/*': runUiChecks
};
