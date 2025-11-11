// VC++, GCC 동작 확인


using   u32     = unsigned int;


#if _MSC_VER

#include <intrin.h>


inline  auto    __builtin_ctz( const u32 bits )

{

    u32     bit = 0;

    _BitScanReverse( (unsigned long*)&bit, bits );

    return  bit;

}

//  BSF BSR 명령을 꼼수로 사용

#endif


inline  auto    where_zero_byte( const __m128i* s )

{

    const   __m128i zero = _mm_set1_epi8( 0 );

    //  8bit char 값 0를 128비트 레지스터 ( 16개 ) 에 채움

    return  _mm_movemask_epi8( _mm_cmpeq_epi8( _mm_loadu_si128( s ), zero ) );

    //  128비트 포인터의 값을 불러와( 16바이트로 정렬되지 않아도 됨 ),

    //  각 바이트 == 0 ? 0xFF : 0x00 로 채우고

    //  인자의 MSB ( 부호비트 ) 를 LSB 부터 차례로 16비트에 채운뒤 남는 앞 16비트는 0으로 채움.

    //  한 마디로 MSB 만 차곡차곡 모아서 16비트를 만듦.  

}

    

auto            strlen( const char* s )

{

    auto*   p = (const __m128i*)s;

    u32     finder;

    while( !( finder = where_zero_byte( p++ ) ) );

    return  (const char*)p - s - sizeof *p + __builtin_ctz( finder );

}


