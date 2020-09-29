package Isuumo::Web;
use v5.32;
use warnings;
use utf8;

use feature qw(isa);
no warnings qw(experimental::isa);

use Kossy;

use DBIx::Sunny;
use File::Spec;
use HTTP::Status qw/:constants/;
use Log::Minimal;
use JSON::MaybeXS;
use Cpanel::JSON::XS::Type;
use Text::CSV_XS;

local $Log::Minimal::LOG_LEVEL = "DEBUG";

my $MYSQL_CONNECTION_DATA = {
    host     => $ENV{MYSQL_HOST}     // '127.0.0.1',
    port     => $ENV{MYSQL_PORT}     // '3306',
    user     => $ENV{MYSQL_USER}     // 'isucon',
    dbname   => $ENV{MYSQL_DBNAME}   // 'isuumo',
    password => $ENV{MYSQL_PASS}     // 'isucon',
};

my $CHAIR_SEARCH_CONDITION;
my $ESTATE_SEARCH_CONDITION;

use constant LIMIT => 20;
use constant NAZOTTE_LIMIT => 50;

use constant InitializeResponse => {
    language => JSON_TYPE_STRING
};

use constant Chair => {
    id          => JSON_TYPE_INT,
    name        => JSON_TYPE_STRING,
    description => JSON_TYPE_STRING,
    thumbnail   => JSON_TYPE_STRING,
    price       => JSON_TYPE_INT,
    height      => JSON_TYPE_INT,
    width       => JSON_TYPE_INT,
    depth       => JSON_TYPE_INT,
    color       => JSON_TYPE_STRING,
    features    => JSON_TYPE_STRING,
    kind        => JSON_TYPE_STRING,
    popularity  => undef,
    stock       => undef,
};

use constant ChairSearchResponse => {
    count  => JSON_TYPE_INT,
    chairs => json_type_arrayof(Chair),
};

use constant ChairListResponse => {
    chairs => json_type_arrayof(Chair),
};

# Estate 物件
use constant Estate => {
    id          => JSON_TYPE_INT,
    thumbnail   => JSON_TYPE_STRING,
    name        => JSON_TYPE_STRING,
    description => JSON_TYPE_STRING,
    latitude    => JSON_TYPE_FLOAT,
    longitude   => JSON_TYPE_FLOAT,
    address     => JSON_TYPE_STRING,
    rent        => JSON_TYPE_INT,
    doorHeight  => JSON_TYPE_INT,
    doorWidth   => JSON_TYPE_INT,
    features    => JSON_TYPE_STRING,
    popularity  => undef,
};

use constant EstateSearchResponse => {
    count   => JSON_TYPE_INT,
    estates => json_type_arrayof(Estate),
};

use constant EstateListResponse => {
    estates => json_type_arrayof(Estate),
};

use constant Coordinate => {
    latitude    => JSON_TYPE_FLOAT,
    longitude   => JSON_TYPE_FLOAT,
};

use constant Coordinates => {
    coordinates => json_type_arrayof(Coordinate),
};

use constant Range => {
    id  => JSON_TYPE_INT,
    min => JSON_TYPE_INT,
    max => JSON_TYPE_INT,
};

use constant RangeCondition => {
    prefix => JSON_TYPE_STRING,
    suffix => JSON_TYPE_STRING,
    ranges => json_type_arrayof(Range),
};

use constant ListCondition => {
    list => json_type_arrayof(JSON_TYPE_STRING),
};

use constant EstateSearchCondition => {
    doorWidth  => RangeCondition,
    doorHeight => RangeCondition,
    rent       => RangeCondition,
    feature    => ListCondition,
};

use constant ChairSearchCondition => {
    width   => RangeCondition,
    height  => RangeCondition,
    depth   => RangeCondition,
    price   => RangeCondition,
    color   => ListCondition,
    feature => ListCondition,
    kind    => ListCondition,
};

my $_JSON = JSON::MaybeXS->new()->allow_blessed(1)->convert_blessed(1)->ascii(1);

$CHAIR_SEARCH_CONDITION = do {
    my $file = File::Spec->catfile("..", "fixture", "chair_condition.json");
    open my $fh, "<:encoding(utf8)", $file or die "cannot open $file";
    my $json = do { local $/; <$fh> };
    $_JSON->decode($json);
};

$ESTATE_SEARCH_CONDITION = do {
    my $file = File::Spec->catfile("..", "fixture", "estate_condition.json");
    open my $fh, "<:encoding(utf8)", $file or die "cannot open $file";
    my $json = do { local $/; <$fh> };
    $_JSON->decode($json);
};



# Initialize
post '/initialize' => \&initialize;

# Chair Handler
get  '/api/chair/{id:\d+}'         => \&get_chair_detail;
post '/api/chair'                  => \&post_chair;
get  '/api/chair/search'           => \&search_chairs;
get  '/api/chair/low_priced'       => \&get_low_priced_chair;
get  '/api/chair/search/condition' => \&get_chair_search_condition;
post '/api/chair/buy/{id:\d+}'     => [qw/allow_json_request/] => \&buy_chair;

# Estate Handler
get  '/api/estate/{id:\d+}'             => \&get_estate_detail;
post '/api/estate'                      => \&post_estate;
get  '/api/estate/search'               => \&search_estates;
get  '/api/estate/low_priced'           => \&get_low_priced_estate;
post '/api/estate/req_doc/{id:\d+}'     => [qw/allow_json_request/] => \&post_estate_request_document;
post '/api/estate/nazotte'              => [qw/allow_json_request/] => \&search_estate_nazotte;
get  '/api/estate/search/condition'     => \&get_estate_search_condition;
get  '/api/recommended_estate/{id:\d+}' => \&search_recommended_estate_with_chair;




sub initialize {
    my ($self, $c)  = @_;

    my $sql_dir = File::Spec->catfile($self->root_dir, '..', 'mysql', 'db');
    my @path = (
        File::Spec->catfile($sql_dir, '0_Schema.sql'),
        File::Spec->catfile($sql_dir, '1_DummyEstateData.sql'),
        File::Spec->catfile($sql_dir, '2_DummyChairData.sql'),
    );
    for my $p (@path) {
        my $cmd = sprintf("mysql -h %s -u %s -p%s -P %s %s < %s",
            $MYSQL_CONNECTION_DATA->{host},
            $MYSQL_CONNECTION_DATA->{user},
            $MYSQL_CONNECTION_DATA->{password},
            $MYSQL_CONNECTION_DATA->{port},
            $MYSQL_CONNECTION_DATA->{dbname},
            $p
        );
        if (my $e = system($cmd)) {
            infof('Initialize script error : %s , %s', $e, $!);
            return $self->res_no_content($c, HTTP_INTERNAL_SERVER_ERROR);
        }
    }

    $self->res_json($c, {
        language => "perl",
    }, InitializeResponse);
};

sub get_chair_detail {
    my ($self, $c)  = @_;

    my $chair_id = $c->args->{id};
    my $query = 'SELECT * FROM chair WHERE id = ?';
    my $chair = $self->dbh->select_row($query, $chair_id);

    if (!$chair) {
        infof("requested id's chair not found : %s", $chair_id);
        return $self->res_no_content($c, HTTP_NOT_FOUND)
    }

    if ($chair->{stock} <= 0) {
        infof("requested id's chair is sold out : %s", $chair_id);
        return $self->res_no_content($c, HTTP_NOT_FOUND)
    }

    return $self->res_json($c, {
        id          => $chair->{id},
        name        => $chair->{name},
        description => $chair->{description},
        thumbnail   => $chair->{thumbnail},
        price       => $chair->{price},
        height      => $chair->{height},
        width       => $chair->{width},
        depth       => $chair->{depth},
        color       => $chair->{color},
        features    => $chair->{features},
        kind        => $chair->{kind},
    }, Chair)
};

sub post_chair {
    my ($self, $c)  = @_;

    my $file = $c->req->uploads->{'chairs'};
    if (!$file) {
        critf("failed to get form file");
        return $self->res_no_content($c, HTTP_BAD_REQUEST);
    }

    my $fh;
    if (!open $fh, "<:encoding(utf8)", $file->path) {
        critf("failed to open form file %s. %s", $file->path, $!);
        return $self->res_no_content($c, HTTP_INTERNAL_SERVER_ERROR);
    }

    my $csv = Text::CSV_XS->new({binary => 1});
    my $dbh = $self->dbh;
    my $txn = $dbh->txn_scope;

    eval {
        while (my $row = $csv->getline($fh)) {
            my ($id, $name, $description, $thumbnail, $price, $height, $width, $depth, $color, $features, $kind, $popularity, $stock) = $row->@*;
            $dbh->query(
                "INSERT INTO chair(id, name, description, thumbnail, price, height, width, depth, color, features, kind, popularity, stock) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)",
                $id, $name, $description, $thumbnail, $price, $height, $width, $depth, $color, $features, $kind, $popularity, $stock
            );
        }
        $txn->commit;
    };
    if ($@) {
        $txn->rollback;
        critf("failed to commit txn: %s", $@);
        return $self->res_no_content($c, HTTP_INTERNAL_SERVER_ERROR);
    }

    $fh->close;
    return $self->res_no_content($c, HTTP_CREATED);
};

sub search_chairs {
    my ($self, $c)  = @_;

    my @conditions;
    my @params;

    if (exists $c->req->parameters->{priceRangeId}) {
        my $price_range_id = $c->req->parameters->{priceRangeId};
        my ($chair_price, $err) = get_range($CHAIR_SEARCH_CONDITION->{price}, $price_range_id);
        if ($err) {
            infof("priceRangeID invalid, %s : %s", $price_range_id, $err);
            return $self->res_no_content($c, HTTP_BAD_REQUEST);
        }
        if ($chair_price->{min} != -1) {
            push @conditions => "price >= ?";
            push @params => $chair_price->{min};
        }
        if ($chair_price->{max} != -1) {
            push @conditions => "price < ?";
            push @params => $chair_price->{max};
        }
    }

    if (exists $c->req->parameters->{heightRangeId}) {
        my $height_range_id = $c->req->parameters->{heightRangeId};
        my ($chair_height, $err) = get_range($CHAIR_SEARCH_CONDITION->{height}, $height_range_id);
        if ($err) {
            infof("heightRangeID invalid, %s : %s", $height_range_id, $err);
            return $self->res_no_content($c, HTTP_BAD_REQUEST);
        }
        if ($chair_height->{min} != -1) {
            push @conditions => "height >= ?";
            push @params => $chair_height->{min};
        }
        if ($chair_height->{max} != -1) {
            push @conditions => "height < ?";
            push @params => $chair_height->{max};
        }
    }

    if (exists $c->req->parameters->{widthRangeId}) {
        my $width_range_id = $c->req->parameters->{widthRangeId};
        my ($chair_width, $err) = get_range($CHAIR_SEARCH_CONDITION->{width}, $width_range_id);
        if ($err) {
            infof("widthRangeID invalid, %s : %s", $width_range_id, $err);
            return $self->res_no_content($c, HTTP_BAD_REQUEST);
        }
        if ($chair_width->{min} != -1) {
            push @conditions => "width >= ?";
            push @params => $chair_width->{min};
        }
        if ($chair_width->{max} != -1) {
            push @conditions => "width < ?";
            push @params => $chair_width->{max};
        }
    }

    if (exists $c->req->parameters->{depthRangeId}) {
        my $depth_range_id = $c->req->parameters->{depthRangeId};
        my ($chair_depth, $err) = get_range($CHAIR_SEARCH_CONDITION->{depth}, $depth_range_id);
        if ($err) {
            infof("depthRangeID invalid, %s : %s", $depth_range_id, $err);
            return $self->res_no_content($c, HTTP_BAD_REQUEST);
        }
        if ($chair_depth->{min} != -1) {
            push @conditions => "depth >= ?";
            push @params => $chair_depth->{min};
        }
        if ($chair_depth->{max} != -1) {
            push @conditions => "depth < ?";
            push @params => $chair_depth->{max};
        }
    }

    if (my $kind = $c->req->parameters->get('kind')) {
        push @conditions => "kind = ?";
        push @params => $kind;
    }

    if (my $color = $c->req->parameters->get('color')) {
        push @conditions => "color = ?";
        push @params => $color;
    }

    if (my $features = $c->req->parameters->get('features')) {
        for my $f (split /,/, $features) {
            push @conditions => "features LIKE CONCAT('%', ?, '%')";
            push @params => $f;
        }
    }

    if (@conditions == 0) {
        infof("Search condition not found");
        return $self->res_no_content($c, HTTP_BAD_REQUEST);
    }

    push @conditions => "stock > 0";

    my $page = $c->req->parameters->get('page');
    if ($page !~ /^\d+$/) {
        infof("Invalid format page parameter : %s", $page);
        return $self->res_no_content($c, HTTP_BAD_REQUEST);
    }

    my $per_page = $c->req->parameters->get('perPage');
    if ($per_page !~ /^\d+$/) {
        infof("Invalid format per_page parameter : %s", $per_page);
        return $self->res_no_content($c, HTTP_BAD_REQUEST);
    }

    my $search_query = "SELECT * FROM chair WHERE ";
    my $count_query = "SELECT COUNT(*) FROM chair WHERE ";
    my $search_condition = join " AND ", @conditions;
    my $limit_offset = " ORDER BY popularity DESC, id ASC LIMIT ? OFFSET ?";

    my $dbh = $self->dbh;

    my $count = $dbh->select_one($count_query . $search_condition, @params);

    push @params => $per_page, $page * $per_page;
    my $chairs = $dbh->select_all($search_query . $search_condition . $limit_offset, @params);

    return $self->res_json($c, {
        count  => $count,
        chairs => [map {
            +{
                id          => $_->{id},
                name        => $_->{name},
                description => $_->{description},
                thumbnail   => $_->{thumbnail},
                price       => $_->{price},
                height      => $_->{height},
                width       => $_->{width},
                depth       => $_->{depth},
                color       => $_->{color},
                features    => $_->{features},
                kind        => $_->{kind},
            }
        } $chairs->@* ],
    }, ChairSearchResponse);
};

sub buy_chair {
    my ($self, $c) = @_;

    my $email = $c->req->body_parameters->{email};
    if (!$email) {
        infof("post buy chair failed : email not found in request body");
        return $self->res_no_content($c, HTTP_BAD_REQUEST);
    }

    my $chair_id = $c->args->{id};
    my $dbh = $self->dbh;
    my $txn = $dbh->txn_scope;

    eval {
        my $chair = $dbh->select_row("SELECT * FROM chair WHERE id = ? AND stock > 0 FOR UPDATE", $chair_id);
        if (!$chair) {
            infof("buyChair chair id \"%s\" not found", $chair_id);
            die $self->res_no_content($c, HTTP_NOT_FOUND);
        }
        $dbh->query("UPDATE chair SET stock = stock - 1 WHERE id = ?", $chair_id);
        $txn->commit;
    };
    if ($@) {
        $txn->rollback;
        return $@ if $@ isa Plack::Response;

        critf("transaction commit error : %s", $@);
        return $self->res_no_content($c, HTTP_INTERNAL_SERVER_ERROR);
    }

    return $self->res_no_content($c, HTTP_OK);
};

sub get_chair_search_condition {
    my ($self, $c) = @_;
    return $self->res_json($c, $CHAIR_SEARCH_CONDITION, ChairSearchCondition);
};

sub get_low_priced_chair {
    my ($self, $c) = @_;

    my $query = "SELECT * FROM chair WHERE stock > 0 ORDER BY price ASC, id ASC LIMIT ?";
    my $chairs = $self->dbh->select_all($query, LIMIT);
    if ($chairs->@* == 0) {
        critf("getLowPricedChair not found");
    }

    return $self->res_json($c, {
        chairs => [map {
            +{
                id          => $_->{id},
                name        => $_->{name},
                description => $_->{description},
                thumbnail   => $_->{thumbnail},
                price       => $_->{price},
                height      => $_->{height},
                width       => $_->{width},
                depth       => $_->{depth},
                color       => $_->{color},
                features    => $_->{features},
                kind        => $_->{kind},
            }
        } $chairs->@* ],
    }, ChairListResponse);
};

sub get_estate_detail {
    my ($self, $c)  = @_;

    my $estate_id = $c->args->{id};
    my $query = 'SELECT * FROM estate WHERE id = ?';
    my $estate = $self->dbh->select_row($query, $estate_id);

    if (!$estate) {
        infof("getEstateDetail estate id not found : %s", $estate_id);
        return $self->res_no_content($c, HTTP_NOT_FOUND)
    }

    return $self->res_json($c, {
        id          => $estate->{id},
        thumbnail   => $estate->{thumbnail},
        name        => $estate->{name},
        description => $estate->{description},
        latitude    => $estate->{latitude},
        longitude   => $estate->{longitude},
        address     => $estate->{address},
        rent        => $estate->{rent},
        doorHeight  => $estate->{door_height},
        doorWidth   => $estate->{door_width},
        features    => $estate->{features},
    }, Estate)
};

sub post_estate {
    my ($self, $c)  = @_;

    my $file = $c->req->uploads->{'estates'};
    if (!$file) {
        critf("failed to get form file");
        return $self->res_no_content($c, HTTP_BAD_REQUEST);
    }

    my $fh;
    if (!open $fh, "<:encoding(utf8)", $file->path) {
        critf("failed to open form file %s. %s", $file->path, $!);
        return $self->res_no_content($c, HTTP_INTERNAL_SERVER_ERROR);
    }

    my $csv = Text::CSV_XS->new({binary => 1});
    my $dbh = $self->dbh;
    my $txn = $dbh->txn_scope;

    eval {
        while (my $row = $csv->getline($fh)) {
            my ($id, $name, $description, $thumbnail, $address, $latitude, $longitude, $rent, $door_height, $door_width, $features, $popularity) = $row->@*;
            $dbh->query(
                "INSERT INTO estate(id, name, description, thumbnail, address, latitude, longitude, rent, door_height, door_width, features, popularity) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)",
                $id, $name, $description, $thumbnail, $address, $latitude, $longitude, $rent, $door_height, $door_width, $features, $popularity
            );
        }
        $txn->commit;
    };
    if ($@) {
        $txn->rollback;
        critf("failed to commit txn: %s", $@);
        return $self->res_no_content($c, HTTP_INTERNAL_SERVER_ERROR);
    }

    $fh->close;
    return $self->res_no_content($c, HTTP_CREATED);
};

sub search_estates {
    my ($self, $c)  = @_;

    my @conditions;
    my @params;

    if (exists $c->req->parameters->{doorHeightRangeId}) {
        my $door_height_range_id = $c->req->parameters->{doorHeightRangeId};
        my ($door_height, $err) = get_range($ESTATE_SEARCH_CONDITION->{doorHeight}, $door_height_range_id);
        if ($err) {
            infof("doorHeightRangeID invalid, %s : %s", $door_height_range_id, $err);
            return $self->res_no_content($c, HTTP_BAD_REQUEST);
        }
        if ($door_height->{min} != -1) {
            push @conditions => "door_height >= ?";
            push @params => $door_height->{min};
        }
        if ($door_height->{max} != -1) {
            push @conditions => "door_height < ?";
            push @params => $door_height->{max};
        }
    }

    if (exists $c->req->parameters->{doorWidthRangeId}) {
        my $door_width_range_id = $c->req->parameters->{doorWidthRangeId};
        my ($door_width, $err) = get_range($ESTATE_SEARCH_CONDITION->{doorWidth}, $door_width_range_id);
        if ($err) {
            infof("doorWidthRangeID invalid, %s : %s", $door_width_range_id, $err);
            return $self->res_no_content($c, HTTP_BAD_REQUEST);
        }
        if ($door_width->{min} != -1) {
            push @conditions => "door_width >= ?";
            push @params => $door_width->{min};
        }
        if ($door_width->{max} != -1) {
            push @conditions => "door_width < ?";
            push @params => $door_width->{max};
        }
    }

    if (exists $c->req->parameters->{rentRangeId}) {
        my $rent_range_id = $c->req->parameters->{rentRangeId};
        my ($estate_rent, $err) = get_range($ESTATE_SEARCH_CONDITION->{rent}, $rent_range_id);
        if ($err) {
            infof("rentRangeID invalid, %s : %s", $rent_range_id, $err);
            return $self->res_no_content($c, HTTP_BAD_REQUEST);
        }
        if ($estate_rent->{min} != -1) {
            push @conditions => "rent >= ?";
            push @params => $estate_rent->{min};
        }
        if ($estate_rent->{max} != -1) {
            push @conditions => "rent < ?";
            push @params => $estate_rent->{max};
        }
    }

    if (my $features = $c->req->parameters->get('features')) {
        for my $f (split /,/, $features) {
            push @conditions => "features LIKE CONCAT('%', ?, '%')";
            push @params => $f;
        }
    }

    if (@conditions == 0) {
        infof("searchEstates search condition not found");
        return $self->res_no_content($c, HTTP_BAD_REQUEST);
    }

    my $page = $c->req->parameters->get('page');
    if ($page !~ /^\d+$/) {
        infof("Invalid format page parameter : %s", $page);
        return $self->res_no_content($c, HTTP_BAD_REQUEST);
    }

    my $per_page = $c->req->parameters->get('perPage');
    if ($per_page !~ /^\d+$/) {
        infof("Invalid format per_page parameter : %s", $per_page);
        return $self->res_no_content($c, HTTP_BAD_REQUEST);
    }

    my $search_query = "SELECT * FROM estate WHERE ";
    my $count_query = "SELECT COUNT(*) FROM estate WHERE ";
    my $search_condition = join " AND ", @conditions;
    my $limit_offset = " ORDER BY popularity DESC, id ASC LIMIT ? OFFSET ?";

    my $dbh = $self->dbh;

    my $count = $dbh->select_one($count_query . $search_condition, @params);

    push @params => $per_page, $page * $per_page;
    my $estates = $dbh->select_all($search_query . $search_condition . $limit_offset, @params);

    return $self->res_json($c, {
        count   => $count,
        estates => [map {
            +{
                id          => $_->{id},
                thumbnail   => $_->{thumbnail},
                name        => $_->{name},
                description => $_->{description},
                latitude    => $_->{latitude},
                longitude   => $_->{longitude},
                address     => $_->{address},
                rent        => $_->{rent},
                doorHeight  => $_->{door_height},
                doorWidth   => $_->{door_width},
                features    => $_->{features},
            }
        } $estates->@* ],
    }, EstateSearchResponse);
};

sub get_low_priced_estate {
    my ($self, $c) = @_;

    my $query = "SELECT * FROM estate ORDER BY rent ASC, id ASC LIMIT ?";
    my $estates = $self->dbh->select_all($query, LIMIT);
    if ($estates->@* == 0) {
        critf("getLowPricedEstate not found");
    }

    return $self->res_json($c, {
        estates => [map {
            +{
                id          => $_->{id},
                thumbnail   => $_->{thumbnail},
                name        => $_->{name},
                description => $_->{description},
                latitude    => $_->{latitude},
                longitude   => $_->{longitude},
                address     => $_->{address},
                rent        => $_->{rent},
                doorHeight  => $_->{door_height},
                doorWidth   => $_->{door_width},
                features    => $_->{features},
            }
        } $estates->@* ],
    }, EstateListResponse);
};

sub search_recommended_estate_with_chair {
    my ($self, $c) = @_;

    my $chair_id = $c->args->{id};
    my $chair_query = "SELECT * FROM chair WHERE id = ?";
    my $chair = $self->dbh->select_row($chair_query, $chair_id);
    if (!$chair) {
        infof("Requested chair id \"%v\" not found", $chair_id);
        return $self->res_no_content(HTTP_BAD_REQUEST);
    }

    my $w = $chair->{width};
    my $h = $chair->{height};
    my $d = $chair->{depth};
    my $query = "SELECT * FROM estate WHERE (door_width >= ? AND door_height >= ?) OR (door_width >= ? AND door_height >= ?) OR (door_width >= ? AND door_height >= ?) OR (door_width >= ? AND door_height >= ?) OR (door_width >= ? AND door_height >= ?) OR (door_width >= ? AND door_height >= ?) ORDER BY popularity DESC, id ASC LIMIT ?";
    my $estates = $self->dbh->select_all($query,
        $w, $h, $w, $d, $h, $w, $h, $d, $d, $w, $d, $h, LIMIT
    );

    return $self->res_json($c, {
        estates => [map {
            +{
                id          => $_->{id},
                thumbnail   => $_->{thumbnail},
                name        => $_->{name},
                description => $_->{description},
                latitude    => $_->{latitude},
                longitude   => $_->{longitude},
                address     => $_->{address},
                rent        => $_->{rent},
                doorHeight  => $_->{door_height},
                doorWidth   => $_->{door_width},
                features    => $_->{features},
            }
        } $estates->@* ],
    }, EstateListResponse);
};

sub search_estate_nazotte {
    my ($self, $c) = @_;

    my $coordinates = $c->req->body_parameters_raw->{coordinates};
    if (!$coordinates->@*) {
        return $self->res_no_content($c, HTTP_BAD_REQUEST);
    }

    my $box = get_bounding_box($coordinates);
    my $query = 'SELECT * FROM estate WHERE latitude <= ? AND latitude >= ? AND longitude <= ? AND longitude >= ? ORDER BY popularity DESC, id ASC';

    my $estates_in_bounding_box = $self->dbh->select_all($query,
        $box->{bottom_right_corner}{latitude}, $box->{top_left_corner}{latitude}, $box->{bottom_right_corner}{longitude}, $box->{top_left_corner}{longitude});

    if ($estates_in_bounding_box->@* == 0) {
        infof("select * from estate where latitude ... not found");
    }

    my $estates_in_polygon = [];
    for my $estate ($estates_in_bounding_box->@*) {
        my $point = sprintf("'POINT(%f %f)'", $estate->{latitude}, $estate->{longitude});
        my $query = sprintf("SELECT * FROM estate WHERE id = ? AND ST_Contains(ST_PolygonFromText(%s), ST_GeomFromText(%s))", coordinates_to_text($coordinates), $point);

        my $validated_estate = $self->dbh->select_row($query, $estate->{id});
        if ($validated_estate) {
            push $estates_in_polygon->@* => $validated_estate;
        }
    }

    my $estates;
    if ($estates_in_polygon->@* > NAZOTTE_LIMIT) {
        $estates = [$estates_in_polygon->@[0.. NAZOTTE_LIMIT - 1]];
    }
    else {
        $estates = $estates_in_polygon;
    }
    my $count = $estates->@*;

    return $self->res_json($c, {
        count   => $count,
        estates => [map {
            +{
                id          => $_->{id},
                thumbnail   => $_->{thumbnail},
                name        => $_->{name},
                description => $_->{description},
                latitude    => $_->{latitude},
                longitude   => $_->{longitude},
                address     => $_->{address},
                rent        => $_->{rent},
                doorHeight  => $_->{door_height},
                doorWidth   => $_->{door_width},
                features    => $_->{features},
            }
        } $estates->@*],
    }, EstateSearchResponse);
};

sub post_estate_request_document {
    my ($self, $c) = @_;

    my $email = $c->req->body_parameters->{email};
    if (!$email) {
        infof("post request document failed : email not found in request body");
        return $self->res_no_content($c, HTTP_BAD_REQUEST);
    }

    my $estate_id = $c->args->{id};
    my $query = 'SELECT * FROM estate WHERE id = ?';
    my $estate = $self->dbh->select_row($query, $estate_id);
    if (!$estate) {
        return $self->res_no_content($c, HTTP_NOT_FOUND);
    }

    return $self->res_no_content($c, HTTP_OK);
};

sub get_estate_search_condition {
    my ($self, $c) = @_;
    return $self->res_json($c, $ESTATE_SEARCH_CONDITION, EstateSearchCondition);
};




sub get_bounding_box {
    my ($coordinates) = @_;
    my $bounding_box = {
        top_left_corner => {
            latitude => $coordinates->[0]->{latitude},
            longitude => $coordinates->[0]->{longitude},
        },,
        bottom_right_corner => {
            latitude => $coordinates->[0]->{latitude},
            longitude => $coordinates->[0]->{longitude},
        }
    };

    for my $coordinate ($coordinates->@*) {
        if ($bounding_box->{top_left_corner}{latitude} > $coordinate->{latitude}) {
            $bounding_box->{top_left_corner}{latitude} = $coordinate->{latitude}
        }
        if ($bounding_box->{top_left_corner}{longitude} > $coordinate->{longitude}) {
            $bounding_box->{top_left_corner}{longitude} = $coordinate->{longitude}
        }

        if ($bounding_box->{bottom_right_corner}{latitude} < $coordinate->{latitude}) {
            $bounding_box->{bottom_right_corner}{latitude} = $coordinate->{latitude}
        }
        if ($bounding_box->{bottom_right_corner}{longitude} < $coordinate->{longitude}) {
            $bounding_box->{bottom_right_corner}{longitude} = $coordinate->{longitude}
        }
    }
    return $bounding_box;
}

sub coordinates_to_text {
    my ($coordinates) = @_;

    my @points;
    for my $c ($coordinates->@*) {
        push @points => sprintf("%f %f", $c->{latitude}, $c->{longitude})
    }
    return sprintf("'POLYGON((%s))'", join ",", @points);
}

sub get_range {
    my ($range_condition, $range_id) = @_;

    my $ranges = $range_condition->{ranges};

    if ($range_id < 0 || $ranges->@* <= $range_id) {
        return undef, "Unexpected Range ID"
    }
    return $ranges->[$range_id], undef
}


filter 'allow_json_request' => sub {
    my $app = shift;
    return sub {
        my ($self, $c) = @_;
        $c->env->{'kossy.request.parse_json_body'} = 1;
        $app->($self, $c);
    };
};

sub dbh {
    my $self = shift;
    $self->{_dbh} ||= do {
        my ($host, $port, $user, $dbname, $password) = $MYSQL_CONNECTION_DATA->@{qw/host port user dbname password/};
        my $dsn = "dbi:mysql:database=$dbname;host=$host;port=$port";
        DBIx::Sunny->connect($dsn, $user, $password, {
            mysql_enable_utf8mb4 => 1,
            mysql_auto_reconnect => 1,
            Callbacks => {
                connected => sub {
                    my $dbh = shift;
                    # XXX $dbh->do('SET SESSION sql_mode="STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION"');
                    return;
                },
            },
        });
    };
}

# send empty body with status code
sub res_no_content {
    my ($self, $c, $status) = @_;
    $c->res->headers->remove_content_headers;
    $c->res->content_length(0);
    $c->res->code($status);
    $c->res;
}

# render_json with json spec
# XXX: $json_specを指定できるようにKossy::Conection#render_jsonを調整
sub res_json {
    my ($self, $c, $obj, $json_spec) = @_;

    my $body = $_JSON->encode($obj, $json_spec);
    $body = $c->escape_json($body);

    if ( ( $c->req->env->{'HTTP_USER_AGENT'} || '' ) =~ m/Safari/ ) {
        $body = "\xEF\xBB\xBF" . $body;
    }

    $c->res->status( 200 );
    $c->res->content_type('application/json; charset=UTF-8');
    $c->res->header( 'X-Content-Type-Options' => 'nosniff' ); # defense from XSS
    $c->res->body( $body );
    $c->res;
}

1;
