require 'tmpdir'

task :"foursquare-crosspost", ["pkg"] do
  sh "go build -o pkg/foursquare-crosspost ./cmd/foursquare-crosspost"
  Dir.mktmpdir do |dir|
    cp "pkg/foursquare-crosspost", "#{dir}/main"
    sh "zip -j pkg/foursquare-crosspost.zip #{dir}/main"
  end
end

directory "pkg"
