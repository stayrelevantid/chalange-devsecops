# Network compliance checks

control 'ssh-port-closed' do
  impact 1.0
  title 'SSH port should be closed'
  desc 'Port 22 should not be open to public internet.'

  describe port(22) do
    it { should_not be_listening }
  end
end

control 'api-port-listening' do
  impact 0.7
  title 'API port should be listening'
  desc 'SecureBank API must be listening on port 8080.'

  describe port(8080) do
    it { should be_listening }
  end
end