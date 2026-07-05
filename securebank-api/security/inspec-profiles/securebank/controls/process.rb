# Process compliance checks

control 'no-root-process' do
  impact 1.0
  title 'Application should not run as root'
  desc 'SecureBank API binary must not run as root user.'

  describe processes('securebank') do
    its('users') { should_not include 'root' }
  end
end