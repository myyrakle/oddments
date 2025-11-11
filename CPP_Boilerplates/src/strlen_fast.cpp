constexpr       std::size_t  finder = ( std::size_t )0x0101010101010101ULL;

constexpr       std::size_t  masker = ( std::size_t )0x8080808080808080ULL;


constexpr       std::size_t  has_zero_7bit( const std::size_t n )

{

    return  ( n - finder ) & masker;

}


constexpr       std::size_t  has_zero_8bit( const std::size_t n )

{

    return  has_zero_7bit( n ) & ~n;

}


inline  auto    where_zero( const std::size_t* w )

{

    if( ( *w & 0x00000000000000FF ) == 0 )  return  (char*)w;

    if( ( *w & 0x000000000000FF00 ) == 0 )  return  (char*)w + 1;

    if( ( *w & 0x0000000000FF0000 ) == 0 )  return  (char*)w + 2;

    if( sizeof finder == 8 )

    {

    if( ( *w & 0x00000000FF000000 ) == 0 )  return  (char*)w + 3;

    if( ( *w & 0x000000FF00000000 ) == 0 )  return  (char*)w + 4;

    if( ( *w & 0x0000FF0000000000 ) == 0 )  return  (char*)w + 5;

    if( ( *w & 0x00FF000000000000 ) == 0 )  return  (char*)w + 6;

    }

    return  (char*)w + ( sizeof finder - 1 );

}


template< int bits >

auto            strlen_bit( const char* s )

{

    const int  step     = sizeof finder == 4 ? 8 : 4;

    const auto has_zero = bits == 7 ? has_zero_7bit : has_zero_8bit;


    if( has_zero( *(std::size_t*)s ) )

        return  where_zero( (std::size_t*)s ) - s;


    auto    w = (std::size_t*)( (size_t)s & ~size_t( sizeof finder - 1 ) ) + 1;


    while( 1 )

    {

        if( has_zero( w[ 0 ] ) ) return  where_zero( w     ) - s;

        if( has_zero( w[ 1 ] ) ) return  where_zero( w + 1 ) - s;

        if( has_zero( w[ 2 ] ) ) return  where_zero( w + 2 ) - s;

        if( has_zero( w[ 3 ] ) ) return  where_zero( w + 3 ) - s;

        if( sizeof finder == 4 )

        {

        if( has_zero( w[ 4 ] ) ) return  where_zero( w + 4 ) - s;

        if( has_zero( w[ 5 ] ) ) return  where_zero( w + 5 ) - s;

        if( has_zero( w[ 6 ] ) ) return  where_zero( w + 6 ) - s;

        if( has_zero( w[ 7 ] ) ) return  where_zero( w + 7 ) - s;

        }

        w += step;

    }

}


auto            strlen_fast( const char* s )

{

    const   size_t  len = strlen_bit< 7 >( s );

    return  len + strlen_bit< 8 >( s + len );

}


int NNN=348;
