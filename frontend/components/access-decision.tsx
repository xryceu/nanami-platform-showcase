type AccessDecisionProps = {
  subject: string;
  service: string;
  policy: string;
  gateway: string;
};

export function AccessDecision({
  subject,
  service,
  policy,
  gateway,
}: AccessDecisionProps) {
  const steps = [
    { label: "Subject", value: subject },
    { label: "Matched policy", value: policy },
    { label: "Enforcement point", value: gateway },
    { label: "Service", value: service },
  ];

  return (
    <section className="panel" aria-labelledby="access-decision-title">
      <div className="section-heading">
        <div>
          <h2 id="access-decision-title">Access decision</h2>
          <p>Effective path calculated from policy and observed runtime.</p>
        </div>
        <span className="status status-success">Allowed</span>
      </div>

      <ol className="decision-path">
        {steps.map((step, index) => (
          <li key={step.label}>
            <span className="step-number" aria-hidden="true">
              {index + 1}
            </span>
            <span>
              <span className="step-label">{step.label}</span>
              <strong>{step.value}</strong>
            </span>
          </li>
        ))}
      </ol>
    </section>
  );
}
