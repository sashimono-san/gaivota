-- Create function to track when records were updated
create or replace function update_updated_at_column()
returns trigger as $$
begin
  new.updated_at = now();
  return new;
  end;
$$ language plpgsql;

-- Create users table
create table users(
  id serial primary key,
  email varchar(320) unique not null,
  first_name varchar not null,
  last_name varchar not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);

create trigger update_users_updated_at before update on users for each row execute procedure update_updated_at_column();

-- Create portfolios table
create table portfolios(
  id serial primary key,
  user_id int references users(id) not null,
  name varchar(50) not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz,
  unique (user_id, name)
);

create trigger update_portfolios_updated_at before update on portfolios for each row execute procedure update_updated_at_column();

-- Create wallets table
create table wallets(
  id serial primary key,
  user_id int references users(id) not null,
  name varchar(50) not null,
  total_value double precision not null default 0.0,
  address varchar not null,
  location varchar(50) not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz,
  unique (user_id, name)
);

create trigger update_wallets_updated_at before update on wallets for each row execute procedure update_updated_at_column();

-- Create investments table
create table investments(
  id serial primary key,
  portfolio_id int references portfolios(id) not null,
  token varchar(50) not null,
  token_symbol varchar(10),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz,
  unique (portfolio_id, token)
);

create trigger update_investments_updated_at before update on investments for each row execute procedure update_updated_at_column();

-- Create positions table
create table positions(
  id serial primary key,
  investment_id int references investments(id) not null,
  amount double precision not null,
  average_price double precision not null,
  profit double precision not null default 0.0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);

create trigger update_positions_updated_at before update on positions for each row execute procedure update_updated_at_column();

-- Create holdings relationship table
create table holdings(
  id serial primary key,
  position_id int references positions(id) not null,
  wallet_id int references wallets(id) not null,
  amount double precision not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);

create trigger update_holdings_updated_at before update on holdings for each row execute procedure update_updated_at_column();

-- Create orders table
create type order_operations as enum ('sell', 'buy');
create type order_types as enum ('limit', 'market');

create table orders(
  id serial primary key,
  position_id int references positions(id) not null,
  amount double precision not null,
  unit_price double precision not null,
  total_price double precision not null,
  operation order_operations not null,
  type order_types not null,
  exchange varchar(50) not null,
  executed_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);

create trigger update_orders_updated_at before update on orders for each row execute procedure update_updated_at_column();

---- create above / drop below ----

-- Drop users table
drop trigger update_users_updated_at on users;
drop table users;

-- Drop portfolios table
drop trigger update_portfolios_updated_at on portfolios;
drop table portfolios;

-- Drop wallets table
drop trigger update_wallets_updated_at on wallets;
drop table wallets;

-- Drop investments table
drop trigger update_investments_updated_at on investments;
drop table investments;

-- Drop positions table
drop trigger update_positions_updated_at on positions;
drop table positions;

-- Drop holdings table
drop trigger update_holdings_updated_at on holdings;
drop table holdings;

-- Drop orders table
drop trigger update_orders_updated_at on orders;
drop table orders;

-- Drop functions
drop function update_updated_at_column();
