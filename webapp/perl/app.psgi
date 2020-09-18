use FindBin;
use lib "$FindBin::Bin/extlib/lib/perl5";
use lib "$FindBin::Bin/lib";
use File::Basename;
use Plack::Builder;
use Isuumo::Web;

my $root_dir = File::Basename::dirname(__FILE__);

my $app = Isuumo::Web->psgi($root_dir);
builder {
    enable 'ReverseProxy';
    $app;
};

# Hello! secret => 'tagomoris';
