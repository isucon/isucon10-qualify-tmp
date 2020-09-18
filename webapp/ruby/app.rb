require 'sinatra'
require 'mysql2'
require 'mysql2-cs-bind'
require 'csv'

class App < Sinatra::Base
  LIMIT = 20
  NAZOTTE_LIMIT = 50
  CHAIR_SEARCH_CONDITION = JSON.parse(File.read('../fixture/chair_condition.json'), symbolize_names: true)
  ESTATE_SEARCH_CONDITION = JSON.parse(File.read('../fixture/estate_condition.json'), symbolize_names: true)

  configure :development do
    require 'sinatra/reloader'
    register Sinatra::Reloader
  end

  configure do
    enable :logging
  end

  set :add_charset, ['application/json']

  helpers do
    def db_info
      {
        host: ENV.fetch('MYSQL_HOST', '127.0.0.1'),
        port: ENV.fetch('MYSQL_PORT', '3306'),
        username: ENV.fetch('MYSQL_USER', 'isucon'),
        password: ENV.fetch('MYSQL_PASS', 'isucon'),
        database: ENV.fetch('MYSQL_DBNAME', 'isuumo'),
      }
    end

    def db
      Thread.current[:db] ||= Mysql2::Client.new(
        host: db_info[:host],
        port: db_info[:port],
        username: db_info[:username],
        password: db_info[:password],
        database: db_info[:database],
        reconnect: true,
        symbolize_keys: true,
      )
    end

    def transaction(name)
      begin_transaction(name)
      yield(name)
      commit_transaction(name)
    rescue Exception => e
      logger.error "Failed to commit tx: #{e.inspect}"
      rollback_transaction(name)
      raise
    ensure
      ensure_to_abort_transaction(name)
    end

    def begin_transaction(name)
      Thread.current[:db_transaction] ||= {}
      db.query('BEGIN')
      Thread.current[:db_transaction][name] = :open
    end

    def commit_transaction(name)
      Thread.current[:db_transaction] ||= {}
      db.query('COMMIT')
      Thread.current[:db_transaction][name] = :nil
    end

    def rollback_transaction(name)
      Thread.current[:db_transaction] ||= {}
      db.query('ROLLBACK')
      Thread.current[:db_transaction][name] = :nil
    end

    def ensure_to_abort_transaction(name)
      Thread.current[:db_transaction] ||= {}
      if in_transaction?(name)
        logger.warn "Transaction closed implicitly (#{$$}, #{Thread.current.object_id}): #{name}"
        rollback_transaction(name)
      end
    end

    def in_transaction?(name)
      Thread.current[:db_transaction] && Thread.current[:db_transaction][name] == :open
    end

    def camelize_keys_for_estate(estate_hash)
      estate_hash.tap do |e|
        e[:doorHeight] = e.delete(:door_height)
        e[:doorWidth] = e.delete(:door_width)
      end
    end

    def body_json_params
      @body_json_params ||= JSON.parse(request.body.tap(&:rewind).read, symbolize_names: true)
    rescue JSON::ParserError => e
      logger.error "Failed to parse body: #{e.inspect}"
      halt 400
    end
  end

  post '/initialize' do
    sql_dir = Pathname.new('../mysql/db')
    %w[0_Schema.sql 1_DummyEstateData.sql 2_DummyChairData.sql].each do |sql|
      sql_path = sql_dir.join(sql)
      cmd = ['mysql', '-h', db_info[:host], '-u', db_info[:username], "-p#{db_info[:password]}", '-P', db_info[:port], db_info[:database]]
      IO.popen(cmd, 'w') do |io|
        io.puts File.read(sql_path)
        io.close
      end
    end

    { language: 'ruby' }.to_json
  end

  get '/api/chair/low_priced' do
    sql = "SELECT * FROM chair WHERE stock > 0 ORDER BY price ASC, id ASC LIMIT #{LIMIT}" # XXX:
    chairs = db.query(sql).to_a
    { chairs: chairs }.to_json
  end

  get '/api/chair/search' do
    search_queries = []
    query_params = []

    if params[:priceRangeId] && params[:priceRangeId].size > 0
      chair_price = CHAIR_SEARCH_CONDITION[:price][:ranges][params[:priceRangeId].to_i]
      unless chair_price
        logger.error "priceRangeID invalid: #{params[:priceRangeId]}"
        halt 400
      end

      if chair_price[:min] != -1
        search_queries << 'price >= ?'
        query_params << chair_price[:min]
      end

      if chair_price[:max] != -1
        search_queries << 'price < ?'
        query_params << chair_price[:max]
      end
    end

    if params[:heightRangeId] && params[:heightRangeId].size > 0
      chair_height = CHAIR_SEARCH_CONDITION[:height][:ranges][params[:heightRangeId].to_i]
      unless chair_height
        logger.error "heightRangeId invalid: #{params[:heightRangeId]}"
        halt 400
      end

      if chair_height[:min] != -1
        search_queries << 'height >= ?'
        query_params << chair_height[:min]
      end

      if chair_height[:max] != -1
        search_queries << 'height < ?'
        query_params << chair_height[:max]
      end
    end

    if params[:widthRangeId] && params[:widthRangeId].size > 0
      chair_width = CHAIR_SEARCH_CONDITION[:width][:ranges][params[:widthRangeId].to_i]
      unless chair_width
        logger.error "widthRangeId invalid: #{params[:widthRangeId]}"
        halt 400
      end

      if chair_width[:min] != -1
        search_queries << 'width >= ?'
        query_params << chair_width[:min]
      end

      if chair_width[:max] != -1
        search_queries << 'width < ?'
        query_params << chair_width[:max]
      end
    end

    if params[:depthRangeId] && params[:depthRangeId].size > 0
      chair_depth = CHAIR_SEARCH_CONDITION[:depth][:ranges][params[:depthRangeId].to_i]
      unless chair_depth
        logger.error "depthRangeId invalid: #{params[:depthRangeId]}"
        halt 400
      end

      if chair_depth[:min] != -1
        search_queries << 'depth >= ?'
        query_params << chair_depth[:min]
      end

      if chair_depth[:max] != -1
        search_queries << 'depth < ?'
        query_params << chair_depth[:max]
      end
    end

    if params[:kind] && params[:kind].size > 0
      search_queries << 'kind = ?'
      query_params << params[:kind]
    end

    if params[:color] && params[:color].size > 0
      search_queries << 'color = ?'
      query_params << params[:color]
    end

    if params[:features] && params[:features].size > 0
      params[:features].split(',').each do |feature_condition|
        search_queries << "features LIKE CONCAT('%', ?, '%')"
        query_params.push(feature_condition)
      end
    end

    if search_queries.size == 0
      logger.error "Search condition not found"
      halt 400
    end

    search_queries.push('stock > 0')

    page =
      begin
        Integer(params[:page], 10)
      rescue ArgumentError => e
        logger.error "Invalid format page parameter: #{e.inspect}"
        halt 400
      end

    per_page =
      begin
        Integer(params[:perPage], 10)
      rescue ArgumentError => e
        logger.error "Invalid format perPage parameter: #{e.inspect}"
        halt 400
      end

    sqlprefix = 'SELECT * FROM chair WHERE '
    search_condition = search_queries.join(' AND ')
    limit_offset = " ORDER BY popularity DESC, id ASC LIMIT #{per_page} OFFSET #{per_page * page}" # XXX: mysql-cs-bind doesn't support escaping variables for limit and offset
    count_prefix = 'SELECT COUNT(*) as count FROM chair WHERE '

    count = db.xquery("#{count_prefix}#{search_condition}", query_params).first[:count]
    chairs = db.xquery("#{sqlprefix}#{search_condition}#{limit_offset}", query_params).to_a

    { count: count, chairs: chairs }.to_json
  end

  get '/api/chair/:id' do
    id =
      begin
        Integer(params[:id], 10)
      rescue ArgumentError => e
        logger.error "Request parameter \"id\" parse error: #{e.inspect}"
        halt 400
      end

    chair = db.xquery('SELECT * FROM chair WHERE id = ?', id).first
    unless chair
      logger.info "Requested id's chair not found: #{id}"
      halt 404
    end

    if chair[:stock] <= 0
      logger.info "Requested id's chair is sold out: #{id}"
      halt 404
    end

    chair.to_json
  end

  post '/api/chair' do
    if !params[:chairs] || !params[:chairs].respond_to?(:key) || !params[:chairs].key?(:tempfile)
      logger.error 'Failed to get form file'
      halt 400
    end

    transaction('post_api_chair') do
      CSV.parse(params[:chairs][:tempfile].read, skip_blanks: true) do |row|
        sql = 'INSERT INTO chair(id, name, description, thumbnail, price, height, width, depth, color, features, kind, popularity, stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)'
        db.xquery(sql, *row.map(&:to_s))
      end
    end

    status 201
  end

  post '/api/chair/buy/:id' do
    unless body_json_params[:email]
      logger.error 'post buy chair failed: email not found in request body'
      halt 400
    end

    id =
      begin
        Integer(params[:id], 10)
      rescue ArgumentError => e
        logger.error "post buy chair failed: #{e.inspect}"
        halt 400
      end

    transaction('post_api_chair_buy') do |tx_name|
      chair = db.xquery('SELECT * FROM chair WHERE id = ? AND stock > 0 FOR UPDATE', id).first
      unless chair
        rollback_transaction(tx_name) if in_transaction?(tx_name)
        halt 404
      end
      db.xquery('UPDATE chair SET stock = stock - 1 WHERE id = ?', id)
    end

    status 200
  end

  get '/api/chair/search/condition' do
    CHAIR_SEARCH_CONDITION.to_json
  end

  get '/api/estate/low_priced' do
    sql = "SELECT * FROM estate ORDER BY rent ASC, id ASC LIMIT #{LIMIT}" # XXX:
    estates = db.xquery(sql).to_a
    { estates: estates.map { |e| camelize_keys_for_estate(e) } }.to_json
  end

  get '/api/estate/search' do
    search_queries = []
    query_params = []

    if params[:doorHeightRangeId] && params[:doorHeightRangeId].size > 0
      door_height = ESTATE_SEARCH_CONDITION[:doorHeight][:ranges][params[:doorHeightRangeId].to_i]
      unless door_height
        logger.error "doorHeightRangeId invalid: #{params[:doorHeightRangeId]}"
        halt 400
      end

      if door_height[:min] != -1
        search_queries << 'door_height >= ?'
        query_params << door_height[:min]
      end

      if door_height[:max] != -1
        search_queries << 'door_height < ?'
        query_params << door_height[:max]
      end
    end

    if params[:doorWidthRangeId] && params[:doorWidthRangeId].size > 0
      door_width = ESTATE_SEARCH_CONDITION[:doorWidth][:ranges][params[:doorWidthRangeId].to_i]
      unless door_width
        logger.error "doorWidthRangeId invalid: #{params[:doorWidthRangeId]}"
        halt 400
      end

      if door_width[:min] != -1
        search_queries << 'door_width >= ?'
        query_params << door_width[:min]
      end

      if door_width[:max] != -1
        search_queries << 'door_width < ?'
        query_params << door_width[:max]
      end
    end

    if params[:rentRangeId] && params[:rentRangeId].size > 0
      rent = ESTATE_SEARCH_CONDITION[:rent][:ranges][params[:rentRangeId].to_i]
      unless rent
        logger.error "rentRangeId invalid: #{params[:rentRangeId]}"
        halt 400
      end

      if rent[:min] != -1
        search_queries << 'rent >= ?'
        query_params << rent[:min]
      end

      if rent[:max] != -1
        search_queries << 'rent < ?'
        query_params << rent[:max]
      end
    end

    if params[:features] && params[:features].size > 0
      params[:features].split(',').each do |feature_condition|
        search_queries << "features LIKE CONCAT('%', ?, '%')"
        query_params.push(feature_condition)
      end
    end

    if search_queries.size == 0
      logger.error "Search condition not found"
      halt 400
    end

    page =
      begin
        Integer(params[:page], 10)
      rescue ArgumentError => e
        logger.error "Invalid format page parameter: #{e.inspect}"
        halt 400
      end

    per_page =
      begin
        Integer(params[:perPage], 10)
      rescue ArgumentError => e
        logger.error "Invalid format perPage parameter: #{e.inspect}"
        halt 400
      end

    sqlprefix = 'SELECT * FROM estate WHERE '
    search_condition = search_queries.join(' AND ')
    limit_offset = " ORDER BY popularity DESC, id ASC LIMIT #{per_page} OFFSET #{per_page * page}" # XXX:
    count_prefix = 'SELECT COUNT(*) as count FROM estate WHERE '

    count = db.xquery("#{count_prefix}#{search_condition}", query_params).first[:count]
    estates = db.xquery("#{sqlprefix}#{search_condition}#{limit_offset}", query_params).to_a

    { count: count, estates: estates.map { |e| camelize_keys_for_estate(e) } }.to_json
  end

  post '/api/estate/nazotte' do
    coordinates = body_json_params[:coordinates]

    unless coordinates
      logger.error "post search estate nazotte failed: coordinates not found"
      halt 400
    end

    if !coordinates.is_a?(Array) || coordinates.empty?
      logger.error "post search estate nazotte failed: coordinates are empty"
      halt 400
    end

    longitudes = coordinates.map { |c| c[:longitude] }
    latitudes = coordinates.map { |c| c[:latitude] }
    bounding_box = {
      top_left: {
        longitude: longitudes.min,
        latitude: latitudes.min,
      },
      bottom_right: {
        longitude: longitudes.max,
        latitude: latitudes.max,
      },
    }

    sql = 'SELECT * FROM estate WHERE latitude <= ? AND latitude >= ? AND longitude <= ? AND longitude >= ? ORDER BY popularity DESC, id ASC'
    estates = db.xquery(sql, bounding_box[:bottom_right][:latitude], bounding_box[:top_left][:latitude], bounding_box[:bottom_right][:longitude], bounding_box[:top_left][:longitude])

    estates_in_polygon = []
    estates.each do |estate|
      point = "'POINT(%f %f)'" % estate.values_at(:latitude, :longitude)
      coordinates_to_text = "'POLYGON((%s))'" % coordinates.map { |c| '%f %f' % c.values_at(:latitude, :longitude) }.join(',')
      sql = 'SELECT * FROM estate WHERE id = ? AND ST_Contains(ST_PolygonFromText(%s), ST_GeomFromText(%s))' % [coordinates_to_text, point]
      e = db.xquery(sql, estate[:id]).first
      if e
        estates_in_polygon << e
      end
    end

    nazotte_estates = estates_in_polygon.take(NAZOTTE_LIMIT)
    {
      estates: nazotte_estates.map { |e| camelize_keys_for_estate(e) },
      count: nazotte_estates.size,
    }.to_json
  end

  get '/api/estate/:id' do
    id =
      begin
        Integer(params[:id], 10)
      rescue ArgumentError => e
        logger.error "Request parameter \"id\" parse error: #{e.inspect}"
        halt 400
      end

    estate = db.xquery('SELECT * FROM estate WHERE id = ?', id).first
    unless estate
      logger.info "Requested id's estate not found: #{id}"
      halt 404
    end

    camelize_keys_for_estate(estate).to_json
  end

  post '/api/estate' do
    unless params[:estates]
      logger.error 'Failed to get form file'
      halt 400
    end

    transaction('post_api_estate') do
      CSV.parse(params[:estates][:tempfile].read, skip_blanks: true) do |row|
        sql = 'INSERT INTO estate(id, name, description, thumbnail, address, latitude, longitude, rent, door_height, door_width, features, popularity) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)'
        db.xquery(sql, *row.map(&:to_s))
      end
    end

    status 201
  end

  post '/api/estate/req_doc/:id' do
    unless body_json_params[:email]
      logger.error 'post request document failed: email not found in request body'
      halt 400
    end

    id =
      begin
        Integer(params[:id], 10)
      rescue ArgumentError => e
        logger.error "post request document failed: #{e.inspect}"
        halt 400
      end

    estate = db.xquery('SELECT * FROM estate WHERE id = ?', id).first
    unless estate
      logger.error "Requested id's estate not found: #{id}"
      halt 404
    end

    status 200
  end

  get '/api/estate/search/condition' do
    ESTATE_SEARCH_CONDITION.to_json
  end

  get '/api/recommended_estate/:id' do
    id =
      begin
        Integer(params[:id], 10)
      rescue ArgumentError => e
        logger.error "Request parameter \"id\" parse error: #{e.inspect}"
        halt 400
      end

    chair = db.xquery('SELECT * FROM chair WHERE id = ?', id).first
    unless chair
      logger.error "Requested id's chair not found: #{id}"
      halt 404
    end

    w = chair[:width]
    h = chair[:height]
    d = chair[:depth]

    sql = "SELECT * FROM estate WHERE (door_width >= ? AND door_height >= ?) OR (door_width >= ? AND door_height >= ?) OR (door_width >= ? AND door_height >= ?) OR (door_width >= ? AND door_height >= ?) OR (door_width >= ? AND door_height >= ?) OR (door_width >= ? AND door_height >= ?) ORDER BY popularity DESC, id ASC LIMIT #{LIMIT}" # XXX:
    estates = db.xquery(sql, w, h, w, d, h, w, h, d, d, w, d, h).to_a

    { estates: estates.map { |e| camelize_keys_for_estate(e) } }.to_json
  end
end
